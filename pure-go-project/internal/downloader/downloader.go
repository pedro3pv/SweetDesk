package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct {
	modelDir string
}

func NewDownloader(modelDir string) *Downloader {
	return &Downloader{modelDir: modelDir}
}

func (d *Downloader) DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// EnsureClassifierModel baixa o modelo de classificação
func (d *Downloader) EnsureClassifierModel() error {
	modelPath := filepath.Join(d.modelDir, "classifier", "model.onnx")
	metaPath := filepath.Join(d.modelDir, "classifier", "meta.json")

	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Println("Baixando model.onnx...")
		url := "https://huggingface.co/deepghs/anime_real_cls/resolve/main/caformer_s36_v1.3_fixed/model.onnx"
		if err := d.DownloadFile(url, modelPath); err != nil {
			return err
		}
		log.Println("✓ model.onnx baixado")
	}

	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		log.Println("Baixando meta.json...")
		url := "https://huggingface.co/deepghs/anime_real_cls/resolve/main/caformer_s36_v1.3_fixed/meta.json"
		if err := d.DownloadFile(url, metaPath); err != nil {
			return err
		}
		log.Println("✓ meta.json baixado")
	}

	return nil
}

// EnsureRealCUGANModel baixa o modelo RealCUGAN
func (d *Downloader) EnsureRealCUGANModel() error {
	modelPath := filepath.Join(d.modelDir, "realcugan", "realcugan-pro.onnx")

	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Println("Baixando RealCUGAN-pro...")
		// URL de exemplo - ajustar para o modelo correto
		url := "https://huggingface.co/deepghs/imgutils-models/resolve/main/real_esrgan/RealESRGAN_x4plus_anime_6B.onnx"
		if err := d.DownloadFile(url, modelPath); err != nil {
			return err
		}
		log.Println("✓ RealCUGAN-pro baixado")
	}

	return nil
}

// EnsureLSDIRModel baixa o modelo 4xLSDIR
func (d *Downloader) EnsureLSDIRModel() error {
	modelPath := filepath.Join(d.modelDir, "lsdir", "4xLSDIR.onnx")

	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Println("Baixando 4xLSDIR...")
		// URL de exemplo - ajustar para o modelo correto
		url := "https://huggingface.co/Jonny001/deepfake/resolve/main/lsdir_x4.onnx"
		if err := d.DownloadFile(url, modelPath); err != nil {
			return err
		}
		log.Println("✓ 4xLSDIR baixado")
	}

	return nil
}

// EnsureAllModels baixa todos os modelos necessários
func (d *Downloader) EnsureAllModels() error {
	if err := d.EnsureClassifierModel(); err != nil {
		return fmt.Errorf("erro ao baixar classificador: %w", err)
	}
	if err := d.EnsureRealCUGANModel(); err != nil {
		return fmt.Errorf("erro ao baixar RealCUGAN: %w", err)
	}
	if err := d.EnsureLSDIRModel(); err != nil {
		return fmt.Errorf("erro ao baixar LSDIR: %w", err)
	}
	return nil
}
