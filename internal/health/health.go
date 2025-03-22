package health

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/syncx"
)
// defaultHealthManager 是一个全局的合并健康管理器
// defaultHealthManager is global comboHealthManager.
var defaultHealthManager = newComboHealthManager()

type (
	// Probe 表示组件健康状态
	// Probe represents readiness status of a given component.
	Probe interface {
		// MarkReady 设置组件为 ready 状态
		// MarkReady sets a ready state for the endpoint handlers.
		MarkReady()
		// MarkReady 设置组件为 not ready 状态
		// MarkNotReady sets a not ready state for the endpoint handlers.
		MarkNotReady()
		// IsReady 返回组件内部状态
		// IsReady return inner state for the component.
		IsReady() bool
		// Name 返回组件名称
		// Name return probe name identifier
		Name() string
	}

	// healthManager manage app healthy.
	healthManager struct {
		ready syncx.AtomicBool
		name  string
	}
    // comboHealthManager 折叠给定的组件状态，在多线程下安全的反映组件状态。
	// comboHealthManager folds given probes into one, reflects their statuses in a thread-safe way.
	comboHealthManager struct {
		mu     sync.Mutex
		probes []Probe
	}
)

// AddProbe add components probe to global comboHealthManager.
func AddProbe(probe Probe) {
	defaultHealthManager.addProbe(probe)
}

// CreateHttpHandler create health http handler base on given probe.
func CreateHttpHandler(healthResponse string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if defaultHealthManager.IsReady() {
			_, _ = w.Write([]byte(healthResponse))
		} else {
			http.Error(w, "Service Unavailable\n"+defaultHealthManager.verboseInfo(),
				http.StatusServiceUnavailable)
		}
	}
}

// NewHealthManager returns a new healthManager.
func NewHealthManager(name string) Probe {
	return &healthManager{
		name: name,
	}
}

// MarkReady sets a ready state for the endpoint handlers.
func (h *healthManager) MarkReady() {
	h.ready.Set(true)
}

// MarkNotReady sets a not ready state for the endpoint handlers.
func (h *healthManager) MarkNotReady() {
	h.ready.Set(false)
}

// IsReady return inner state for the component.
func (h *healthManager) IsReady() bool {
	return h.ready.True()
}

// Name return probe name identifier
func (h *healthManager) Name() string {
	return h.name
}

func newComboHealthManager() *comboHealthManager {
	return &comboHealthManager{}
}

// MarkReady sets components status to ready.
func (p *comboHealthManager) MarkReady() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, probe := range p.probes {
		probe.MarkReady()
	}
}

// MarkNotReady sets components status to not ready with given error as a cause.
func (p *comboHealthManager) MarkNotReady() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, probe := range p.probes {
		probe.MarkNotReady()
	}
}

// IsReady return composed status of all components.
func (p *comboHealthManager) IsReady() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.probes) == 0 {
		return false
	}

	for _, probe := range p.probes {
		if !probe.IsReady() {
			return false
		}
	}

	return true
}

func (p *comboHealthManager) verboseInfo() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	var info strings.Builder
	for _, probe := range p.probes {
		if probe.IsReady() {
			info.WriteString(fmt.Sprintf("%s is ready\n", probe.Name()))
		} else {
			info.WriteString(fmt.Sprintf("%s is not ready\n", probe.Name()))
		}
	}

	return info.String()
}

// addProbe add components probe to comboHealthManager.
func (p *comboHealthManager) addProbe(probe Probe) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.probes = append(p.probes, probe)
}
