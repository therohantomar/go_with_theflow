package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"

	"go-with-the-flow/internal/cli"
	"go-with-the-flow/internal/generator"
)

func main() {
	cmd := "new"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "new":
		if err := runNew(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "ref":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: flow ref <website-url>")
			os.Exit(2)
		}
		if err := runRef(os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("usage: flow new | flow ref <website-url>")
	}
}

func runRef(referenceURL string) error {
	cfg, err := cli.RunReference(referenceURL)
	if err != nil {
		return err
	}
	ok, err := cli.ConfirmReference(cfg)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("Cancelled.")
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	out, err := generator.Generate(cfg, cwd)
	if err != nil {
		return err
	}
	fmt.Println("\n✔ Generated reference-based Astro project", out)
	fmt.Println("  Review src/content and src/styles/global.css: page text, layout, images, and colors may need manual refinement.")
	fmt.Printf("\nNext:\n  cd %s\n  npm install\n  npm run dev\n", cfg.Slug)
	return nil
}

func runNew() error {
	cfg, err := cli.RunWizard()
	if err != nil {
		return err
	}
	if cfg == nil { // user cancelled at review
		fmt.Println("Cancelled.")
		return nil
	}
	cwd, _ := os.Getwd()
	out, err := generator.Generate(cfg, cwd)
	if err != nil {
		return err
	}
	fmt.Println("\n✔ Generated", out)

	install := false
	_ = huh.NewConfirm().Title("Install dependencies with " + cfg.PackageManager + " now?").Value(&install).Run()
	if install {
		c := exec.Command(cfg.PackageManager, "install")
		c.Dir, c.Stdout, c.Stderr = out, os.Stdout, os.Stderr
		if err := c.Run(); err != nil {
			return err
		}
	}
	fmt.Printf("\nNext:\n  cd %s\n", cfg.Slug)
	if !install {
		fmt.Printf("  %s install\n", cfg.PackageManager)
	}
	run := cfg.PackageManager + " run dev"
	if cfg.PackageManager == "npm" {
		run = "npm run dev"
	}
	fmt.Printf("  %s\n", run)
	return nil
}
