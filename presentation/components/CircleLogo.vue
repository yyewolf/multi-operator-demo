<template>
  <div
    class="circle-image"
    :style="{ width: size, height: size, backgroundImage: `url(${src})` }"
  />
</template>

<script setup>
defineProps({
  src: {
    type: String,
    required: true
  },
  size: {
    type: String,
    default: '100px'
  }
})
</script>

<style scoped>
.circle-image {
  background-size: cover;
  background-position: center;
  border-radius: 50%;
  display: inline-block;
  /* Optional: Add border or shadow */
  border: 3px solid white;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}
</style>
