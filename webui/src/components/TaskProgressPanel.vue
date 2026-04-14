<template>
    <div v-if="visible" class="task-progress-card">
        <div class="task-progress-header">
            <strong>{{ title }}</strong>
            <span class="task-progress-badge">{{ statusLabel }}</span>
        </div>
        <div v-if="summary !== ''" class="task-progress-summary">{{ summary }}</div>
        <div v-if="logEntries.length > 0" class="task-progress-logs">
            <pre>{{ renderedLogs }}</pre>
        </div>
    </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
    visible: { type: Boolean, default: false },
    title: { type: String, default: "任务进度" },
    statusLabel: { type: String, default: "" },
    summary: { type: String, default: "" },
    logEntries: { type: Array, default: () => [] },
});

const renderedLogs = computed(() => props.logEntries.map((entry) => {
    const time = typeof entry?.timestamp === "string" && entry.timestamp.trim() !== "" ? `[${entry.timestamp}] ` : "";
    const message = typeof entry?.message === "string" ? entry.message : "";
    return `${time}${message}`.trimEnd();
}).join("\n"));
</script>
