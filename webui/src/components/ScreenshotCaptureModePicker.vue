<template>
    <div class="variant-picker">
        <button
            class="ghost variant-option"
            :class="{ active: modelValue === 'standard' }"
            :disabled="busy"
            type="button"
            @click="$emit('update:modelValue', 'standard')"
        >
            标准
        </button>
        <button
            class="ghost variant-option"
            :class="{ active: modelValue === 'fast' }"
            :disabled="busy"
            type="button"
            @click="$emit('update:modelValue', 'fast')"
        >
            快速
        </button>
    </div>
</template>

<script setup>
defineProps({
    busy: { type: Boolean, required: true },
    modelValue: { type: String, required: true },
});

defineEmits(["update:modelValue"]);
</script>
