package plugin_wasmi

import (
	"fmt"
)

func RunFile(file []byte, r Runtime) error {

	compiledMod, err := r.R.CompileModule(r.Ctx, file)
	if err != nil {
		return err
	}

	_, err = r.R.InstantiateModule(r.Ctx, compiledMod, r.ModuleConfigs)
	if err != nil {
		return fmt.Errorf("failed to start module: %v", err)
	}

	return nil
}
