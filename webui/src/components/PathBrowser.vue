<template>
    <div class="field">
        <label for="path-search">媒体选择</label>
        <div class="path-picker">
            <div class="browser integrated">
                <div class="browser-toolbar">
                    <div class="browser-current">{{ browserDir || "可用挂载路径" }}</div>
                </div>

                <div class="browser-search">
                    <div class="search-actions">
                        <button
                            class="ghost icon-btn"
                            :disabled="busy || browserLoading || !canNavigateUp"
                            title="上一级"
                            aria-label="上一级"
                            @click="$emit('navigate-up')"
                        >
                            ⬆
                        </button>
                        <button
                            class="ghost icon-btn"
                            :disabled="busy || browserLoading"
                            title="刷新"
                            aria-label="刷新"
                            @click="$emit('refresh')"
                        >
                            ↻
                        </button>
                    </div>
                    <input
                        id="path-search"
                        :value="searchKeyword"
                        type="text"
                        placeholder="模糊搜索当前目录"
                        @input="$emit('update:searchKeyword', $event.target.value)"
                    />
                </div>

                <div v-if="browserError !== ''" class="browser-error">
                    {{ browserError }}
                </div>

                <div class="browser-list">
                    <div v-if="browserLoading" class="browser-row empty">加载中...</div>
                    <div v-else-if="entries.length === 0" class="browser-row empty">当前目录无匹配项</div>
                    <div
                        v-for="entry in entries"
                        :key="entry.path"
                        class="browser-row"
                        :class="{
                            selected: normalizeComparePath(path) === normalizeComparePath(entry.path),
                            directory: entry.isDir,
                            locked: busy || browserLoading,
                        }"
                        @click="$emit('update:path', entry.path)"
                        @dblclick="$emit('open-entry', entry)"
                    >
                        <span class="browser-row-name">{{ entry.isDir ? `📁 ${entry.name}` : `📄 ${entry.name}` }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
defineProps({
    path: { type: String, required: true },
    searchKeyword: { type: String, required: true },
    busy: { type: Boolean, required: true },
    browserDir: { type: String, required: true },
    browserError: { type: String, required: true },
    browserLoading: { type: Boolean, required: true },
    canNavigateUp: { type: Boolean, required: true },
    entries: { type: Array, required: true },
});

defineEmits(["update:path", "update:searchKeyword", "navigate-up", "refresh", "open-entry"]);

const normalizeComparePath = (value) => {
    if (!value) {
        return "";
    }
    if (value === "/" || value === "\\") {
        return "/";
    }
    return value.replace(/\\/g, "/").replace(/\/+$/, "").toLowerCase();
};
</script>
