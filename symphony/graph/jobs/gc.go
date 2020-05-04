package jobs

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// syncServices job syncs the services according to changes
func (m *jobs) garbageCollector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := m.collectProperties(ctx); err != nil {
		m.logger.For(ctx).Error("collect properties", zap.Error(err))
	}
	w.WriteHeader(http.StatusOK)
}

func (m *jobs) collectProperties(ctx context.Context) error {
	client := ent.FromContext(ctx)
	m.logger.For(ctx).Info("running properties garbage collect")
	propertyTypes, err := client.PropertyType.Query().
		Where(propertytype.Deleted(true)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query properties: %w", err)
	}
	for _, pType := range propertyTypes {
		m.logger.For(ctx).Info("deleting property type",
			zap.Int("id", pType.ID),
			zap.String("name", pType.Name))
		count, err := client.Property.Delete().
			Where(property.HasTypeWith(propertytype.ID(pType.ID))).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete properties of type id: %d, %w", pType.ID, err)
		}
		m.logger.For(ctx).Info("deleted properties",
			zap.Int("id", pType.ID),
			zap.Int("count", count))
		err = client.PropertyType.DeleteOne(pType).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete property type: %d, %w", pType.ID, err)
		}
		m.logger.For(ctx).Info("deleted property type",
			zap.Int("id", pType.ID),
			zap.String("name", pType.Name))
	}
	return nil
}
