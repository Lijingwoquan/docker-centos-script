package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const defaultUser = "lijingwoquan"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./init.exe <user_name> <user_password> (can not equal)")
		fmt.Println("Example: ./init.exe liuzihao 123456")
		os.Exit(1)
	}

	userName := os.Args[1]
	password := os.Args[2]

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}

	dockerfilePath := filepath.Join(currentDir, "Dockerfile")

	content, err := readDockerfile(dockerfilePath)
	if err != nil {
		fmt.Println("Error reading Dockerfile:", err)
		os.Exit(1)
	}

	updatedContent := replaceUserAndPassword(content, userName, password)

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "Dockerfile-")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		os.Exit(1)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temporary file when done

	if err = writeDockerfile(tempFile.Name(), updatedContent); err != nil {
		fmt.Println("Error writing updated Dockerfile:", err)
		os.Exit(1)
	}

	fmt.Println("Temporary Dockerfile created successfully.")

	imageName := fmt.Sprintf("centos7-%s:latest", userName)
	err = buildDockerImage(imageName, tempFile.Name())
	if err != nil {
		fmt.Println("Error building Docker image:", err)
		os.Exit(1)
	}

	if err = runDockerImage(imageName); err != nil {
		fmt.Printf("Error running Docker image '%s': %v\n", imageName, err)
		os.Exit(1)
	}
	fmt.Printf("Docker image '%s' built and run successfully.\n", imageName)
}

func readDockerfile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func replaceUserAndPassword(content, userName, password string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "useradd") {
			lines[i] = strings.Replace(line, defaultUser, userName, 1)
		} else if strings.Contains(line, "echo") && strings.Contains(line, "chpasswd") {
			lines[i] = fmt.Sprintf(`    echo "%s:%s" | chpasswd && \`, userName, password)
		} else if strings.Contains(line, "usermod") {
			lines[i] = strings.Replace(line, defaultUser, userName, 1)
		} else if strings.Contains(line, "echo") && strings.Contains(line, "sudoers") {
			lines[i] = fmt.Sprintf(`RUN echo "%s ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers`, userName)
		} else if strings.Contains(line, "USER") {
			lines[i] = fmt.Sprintf("USER %s", userName)
		}
	}
	return strings.Join(lines, "\n")
}

func writeDockerfile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func buildDockerImage(imageName string, dockerfilePath string) error {
	cmd := exec.Command("docker", "build", "-t", imageName, "-f", dockerfilePath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDockerImage(imageName string) error {
	// 移除标签，只使用镜像名称的基本部分作为容器名称
	containerName := strings.Split(imageName, ":")[0]

	// 确保容器名称符合 Docker 的命名规则
	containerName = regexp.MustCompile(`[^a-zA-Z0-9_.-]`).ReplaceAllString(containerName, "")

	// 如果名称以连字符开头，去掉它
	containerName = strings.TrimPrefix(containerName, "-")

	// 如果处理后的名称为空，使用一个默认名称
	if containerName == "" {
		containerName = "my-container"
	}

	cmd := exec.Command("docker", "run", "-it", "--name", containerName, imageName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("正在启动容器，名称: %s，镜像: %s\n", containerName, imageName)

	return cmd.Run()
}
