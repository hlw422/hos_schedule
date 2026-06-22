package worker

import (
	"context"
	"log"
	"time"

	"hos_schedule/internal/model"
	redisutil "hos_schedule/internal/pkg/redis"

	"gorm.io/gorm"
)

type SlotReleaseWorker struct {
	db          *gorm.DB
	slotManager *redisutil.SlotManager
	interval    time.Duration
	expireMin   int
}

func NewSlotReleaseWorker(db *gorm.DB, slotManager *redisutil.SlotManager) *SlotReleaseWorker {
	return &SlotReleaseWorker{
		db:          db,
		slotManager: slotManager,
		interval:    1 * time.Minute,
		expireMin:   15,
	}
}

func (w *SlotReleaseWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Println("Slot release worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Slot release worker stopped")
			return
		case <-ticker.C:
			w.releaseExpiredSlots(ctx)
		}
	}
}

func (w *SlotReleaseWorker) releaseExpiredSlots(ctx context.Context) {
	var appointments []model.Appointment

	err := w.db.Where("status = ? AND created_at < ?",
		"PENDING_PAY",
		time.Now().Add(-time.Duration(w.expireMin)*time.Minute),
	).Find(&appointments).Error

	if err != nil {
		log.Printf("Failed to query expired appointments: %v", err)
		return
	}

	if len(appointments) == 0 {
		return
	}

	log.Printf("Found %d expired appointments to release", len(appointments))

	for _, appt := range appointments {
		result := w.db.Model(&model.Appointment{}).
			Where("id = ? AND status = ?", appt.ID, "PENDING_PAY").
			Update("status", "CANCELLED")

		if result.Error != nil {
			log.Printf("Failed to cancel appointment %d: %v", appt.ID, result.Error)
			continue
		}

		if result.RowsAffected == 0 {
			continue
		}

		if err := w.slotManager.ReleaseSlot(ctx, appt.ScheduleID); err != nil {
			log.Printf("Failed to release Redis slot for schedule %d: %v", appt.ScheduleID, err)
		}

		w.db.Model(&model.Schedule{}).
			Where("id = ?", appt.ScheduleID).
			Updates(map[string]interface{}{
				"remain_count": gorm.Expr("remain_count + 1"),
				"used_count":   gorm.Expr("used_count - 1"),
			})

		log.Printf("Released slot for appointment %d (schedule %d)", appt.ID, appt.ScheduleID)
	}
}
