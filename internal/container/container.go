package container

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"

	"github.com/IslamWalid/tcontainer/internal/namegen"
)

const processesLimit = "10"

// Initialize creates the rootfs fro the container
func Initialize() error {
	_, err := os.Stat("/tmp/rootfs")
	if os.IsNotExist(err) {
		script, err := exec.Command("curl", "https://raw.githubusercontent.com/IslamWalid/tcontainer/master/install.sh").Output()
		if err != nil {
			return err
		}
		err = exec.Command("bash", "-c", string(script)).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// Run starts the parent process of the container
func Run(command string, args []string) error {
	cmd := exec.Command("/proc/self/exe", append([]string{"child", command}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(os.Getuid()),
			Gid: uint32(os.Getgid()),
		},
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings:  []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Getuid(), Size: 1}},
		GidMappings:  []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Getgid(), Size: 1}},
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Child starts a child process to execute given command
func Child(command string, args []string) error {
	err := syscall.Sethostname([]byte(namegen.NameGenerator()))
	if err != nil {
		return err
	}

	err = createCgroup()
	if err != nil {
		return err
	}

	err = syscall.Chroot("/tmp/rootfs")
	if err != nil {
		return err
	}

	err = os.Chdir("/")
	if err != nil {
		return err
	}

	err = syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		return err
	}

	err = syscall.Mount("dev", "dev", "tmpfs", 0, "")
	if err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	// Clean up
	err = syscall.Unmount("/proc", 0)
	if err != nil {
		return err
	}

	err = syscall.Unmount("/dev", 0)
	if err != nil {
		return err
	}

	return nil
}

func createCgroup() error {
	containerPidsDir := "/sys/fs/cgroup/pids"
	os.Mkdir(containerPidsDir, 0755)

	err := os.WriteFile(path.Join(containerPidsDir, "pids.max"), []byte(processesLimit), 0700)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(containerPidsDir, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	if err != nil {
		return err
	}

	return nil
}
