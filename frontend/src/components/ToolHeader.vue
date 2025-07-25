<template>
  <div class="tool-header">
    <h2><i :class="icon"></i> {{ title }}</h2>
    <p>{{ description }}</p>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    required: true
  }
})
</script>

<style scoped>
.tool-header {
  margin-bottom: 2rem;
}

.tool-header h2 {
  font-size: 1.8rem;
  color: #333;
  margin-bottom: 0.5rem;
}

.tool-header p {
  color: #666;
  line-height: 1.6;
}
</style>