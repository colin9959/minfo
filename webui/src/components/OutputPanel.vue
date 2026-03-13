<template>
    <section class="panel output">
        <div class="output-header">
            <h2>输出</h2>
            <div class="output-actions">
                <button class="ghost" @click="$emit('copy')">{{ copyLabel }}</button>
                <button class="ghost" :disabled="busy" @click="$emit('clear')">清空</button>
            </div>
        </div>
        <pre>{{ output }}</pre>
    </section>
</template>

<script setup>
defineProps({
    busy: { type: Boolean, required: true },
    copyLabel: { type: String, required: true },
    output: { type: String, required: true },
});

defineEmits(["copy", "clear"]);
</script>
