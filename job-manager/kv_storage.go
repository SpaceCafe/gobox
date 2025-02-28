package job_manager

var (
	_ IJobHookContext = (*KVStorage)(nil)
)

type KVStorage map[string]any

func (r KVStorage) Get(key string) any {
	return r[key]
}

func (r KVStorage) Set(key string, value any) {
	r[key] = value
}
