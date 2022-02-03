package dbx

import (
	"context"
	"fmt"
)

// Execute Reset Schema
func ExecuteReset() error {
	var tbEmployee = employeeEntity.Tablename()
	var tbTimeData = timeDataEntity.Tablename()
	var ctx = context.Background()
	var err error
	err = ExecSQLs(ctx,
		NewBuilder().
			Addf("DROP TABLE IF EXISTS %s;", tbTimeData).
			String(),
		NewBuilder().
			Addf("DROP TABLE IF EXISTS %s;", tbEmployee).
			String(),
		NewBuilder().
			Addf("CREATE TABLE IF NOT EXISTS %s (", tbEmployee).Addln().
			Addln("  id BIGINT NOT NULL AUTO_INCREMENT,").
			Addln("  user_name VARCHAR(255) NOT NULL UNIQUE,").
			Addln("  display_name VARCHAR(255) NOT NULL,").
			Addln("  PRIMARY KEY (id)").
			Addln(") Engine=InnoDB;").
			String(),
		NewBuilder().
			Addf("CREATE TABLE IF NOT EXISTS %s (", tbTimeData).Addln().
			Addln("  id BIGINT NOT NULL AUTO_INCREMENT,").
			Addln("  parent_id BIGINT NOT NULL,").
			Addln("  name VARCHAR(255) NOT NULL,").
			Addln("  title VARCHAR(255) NOT NULL,").
			Addln("  daily_current BIGINT NOT NULL DEFAULT 0,").
			Addln("  daily_previous BIGINT NOT NULL DEFAULT 0,").
			Addln("  weekly_current BIGINT NOT NULL DEFAULT 0,").
			Addln("  weekly_previous BIGINT NOT NULL DEFAULT 0,").
			Addln("  monthly_current BIGINT NOT NULL DEFAULT 0,").
			Addln("  monthly_previous BIGINT NOT NULL DEFAULT 0,").
			Addln("  PRIMARY KEY (id)").
			Addln(") Engine=InnoDB;").
			String(),
		NewBuilder().
			Addf("ALTER TABLE %s", tbTimeData).Addln().
			Addf("  ADD CONSTRAINT fk_%s_parent_id FOREIGN KEY (parent_id) REFERENCES %s(id);", tbTimeData, tbEmployee).Addln().
			String(),
	)
	if err != nil {
		return err
	}
	fmt.Printf("Schema %s recreated", GetSchema())
	fmt.Println()
	return nil
}
