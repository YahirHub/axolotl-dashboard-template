package feature

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
	"github.com/YahirHub/axolotl-dashboard-template/database/migrations"
)

func TestSQLiteOrmRoundTrip(t *testing.T) {
	migration := &migrations.M20260928000100CreateAppSettingsTable{}

	_ = migration.Down()
	require.NoError(t, migration.Up())
	t.Cleanup(func() {
		_ = migration.Down()
	})

	setting := models.AppSetting{
		Key:   "dashboard.theme",
		Value: "axolotl",
	}
	require.NoError(t, facades.Orm().Query().Create(&setting))

	var stored models.AppSetting
	require.NoError(t, facades.Orm().Query().
		Where("key", setting.Key).
		First(&stored))
	require.Equal(t, setting.Key, stored.Key)
	require.Equal(t, setting.Value, stored.Value)
	require.NotZero(t, stored.ID)

	stored.Value = "axolotl-dark"
	require.NoError(t, facades.Orm().Query().Save(&stored))

	var updated models.AppSetting
	require.NoError(t, facades.Orm().Query().
		Where("key", setting.Key).
		First(&updated))
	require.Equal(t, "axolotl-dark", updated.Value)
}
