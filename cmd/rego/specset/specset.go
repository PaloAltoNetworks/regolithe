package specset

import (
	_ "embed"
	"fmt"
	"os"
	"path"
)

var (
	//go:embed data/Makefile
	fMakeFile []byte
	//go:embed data/custom_validations.nogo
	fCustomValidations []byte
	//go:embed data/custom_validations_test.nogo
	fCustomValidationsTest []byte
	//go:embed data/specs/.regolithe-gen-cmd
	fSpecsRegolitheGenCmd []byte
	//go:embed data/specs/@identifiable.abs
	fSpecsIdentifiable []byte
	//go:embed data/specs/_api.info
	fSpecsAPIInfo []byte
	//go:embed data/specs/_parameter.mapping
	fSpecsParameterMapping []byte
	//go:embed data/specs/_type.mapping
	fSpecsTypeMapping []byte
	//go:embed data/specs/_validation.mapping
	fSpecsValidationMapping []byte
	//go:embed data/specs/object.spec
	fSpecsObject []byte
	//go:embed data/specs/regolithe.ini
	fSpecsRegolitheIni []byte
	//go:embed data/specs/root.spec
	fSpecsRoot []byte
)

// Dump dumps the initial specset into the given path.
func Dump(outPath string) error {

	// root content
	if err := os.MkdirAll(outPath, 0700); err != nil {
		return fmt.Errorf("unable to create initial folder: %w", err)
	}
	if err := os.WriteFile(path.Join(outPath, "Makefile"), fMakeFile, 0700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "custom_validations.go"), fCustomValidations, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "custom_validations_test.go"), fCustomValidationsTest, 0600); err != nil {
		return err
	}

	// specs content
	specsDir := path.Join(outPath, "specs")
	if err := os.MkdirAll(specsDir, 0700); err != nil {
		return fmt.Errorf("unable to create specs folder: %w", err)
	}
	if err := os.WriteFile(path.Join(outPath, "specs", ".regolithe-gen-cmd"), fSpecsRegolitheGenCmd, 0700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "specs", "@identifiable.abs"), fSpecsIdentifiable, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "specs", "_api.info"), fSpecsAPIInfo, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "specs", "_parameter.mapping"), fSpecsParameterMapping, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "specs", "_type.mapping"), fSpecsTypeMapping, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "specs", "_validation.mapping"), fSpecsValidationMapping, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "specs", "object.spec"), fSpecsObject, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(outPath, "specs", "regolithe.ini"), fSpecsRegolitheIni, 0600); err != nil {
		return err
	}
	return os.WriteFile(path.Join(outPath, "specs", "root.spec"), fSpecsRoot, 0600)
}
