package v1

import (
	"net/http"
	"strconv"

	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/models"

	"github.com/labstack/echo/v4"
)

// GetDashboardStats returns global statistics for the admin dashboard
// @Summary Get admin dashboard statistics
// @Description Returns global stats (total tasks, overdue, etc.) and per-user stats. Requires admin permissions.
// @tags admin
// @Accept json
// @Produce json
// @Security JWTKeyAuth
// @Success 200 {object} models.DashboardStats "The dashboard statistics"
// @Failure 403 {object} models.Message "Forbidden"
// @Failure 500 {object} models.Message "Internal error"
// @Router /admin/dashboard [get]
func GetDashboardStats(c echo.Context) error {
	s := db.NewSession()
	defer s.Close()

	stats, err := models.GetDashboardStats(s)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

// GetUserDetail returns detailed statistics and tasks for a specific user
// @Summary Get user details for admin
// @Description Returns stats and tasks for a specific user. Requires admin permissions.
// @tags admin
// @Accept json
// @Produce json
// @Security JWTKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.AdminUserDetail "The user details"
// @Failure 403 {object} models.Message "Forbidden"
// @Failure 500 {object} models.Message "Internal error"
// @Router /admin/users/{id} [get]
func GetUserDetail(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	s := db.NewSession()
	defer s.Close()

	stats, err := models.GetAdminUserDetails(s, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}
