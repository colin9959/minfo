<template>
    <div class="variant-picker">
        <button
            type="button"
            class="ghost variant-option"
            :class="{ active: modelValue === 'png' }"
            :disabled="busy"
            @click="$emit('update:modelValue', 'png')"
        >
            PNG
        </button>
        <button
            type="button"
            class="ghost variant-option"
            :class="{ active: modelValue === 'jpg' }"
            :disabled="busy"
            @click="$emit('update:modelValue', 'jpg')"
        >
            JPG
        </button>
    </div>
</template>

<script setup>
defineProps({
    modelValue: { type: String, required: true },
    busy: { type: Boolean, required: true },
});

defineEmits(["update:modelValue"]);
</script>
