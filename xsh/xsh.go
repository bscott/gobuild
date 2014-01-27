package xsh

import (
	"io"
	"os"
	"os/exec"

	"github.com/codegangsta/inject"
)

func Call(a ...interface{}) error {
	return NewSession().Call(a...)
}

type Dir string

type Session struct {
	inj inject.Injector
}

func NewSession(a ...interface{}) *Session {
	s := &Session{
		inj: inject.New(),
	}
	env := map[string]string{
		"PATH": "/bin:/usr/bin:/usr/local/bin",
	}
	dir := Dir("")
	args := []string{}
	s.inj.Map(env).Map(dir).Map(args).Map(os.Stdout)
	//s.inj.MapTo(os.Stdout, (*io.Writer)(nil))
	//fmt.Println(reflect.ValueOf((*io.Writer)(os.Stdout)).Type())
	for _, v := range a {
		if writer, ok := v.(*io.Writer); ok {
			s.inj.MapTo(writer, (*io.Writer)(nil))
			continue
		}
		s.inj.Map(v)
	}
	return s
}

func (s *Session) Call(a ...interface{}) error {
	for _, v := range a {
		s.inj.Map(v)
	}
	values, err := s.inj.Invoke(invokeExec)
	if err != nil {
		return err
	}
	r := values[0]
	if r.IsNil() {
		return nil
	}
	return r.Interface().(error)
}

func invokeExec(cmd string, args []string, env map[string]string, cwd Dir, output *os.File) error {
	envs := make([]string, 0, len(env))
	for k, v := range env {
		envs = append(envs, k+"="+v)
	}
	c := exec.Command(cmd, args...)
	c.Env = envs
	c.Dir = string(cwd)
	c.Stdout = output
	c.Stderr = output
	return c.Run()
}
