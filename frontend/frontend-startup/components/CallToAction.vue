<template>
  <section class="call-to-action bg-purple-progress pt-64 pb-10">
    <div class="container mx-auto">
      <div class="w-full text-center">
        <h1 class="text-5xl text-white font-semibold">
          Easy way to funding
          <br />
          best idea and innovation
        </h1>
        <button
          class="inline-block bg-orange-button hover:bg-green-button text-white font-semibold px-6 py-4 mt-8 text-lg rounded-full"
          @click="$router.push({ path: '/upload' })"
        >
          Getting Start
        </button>
      </div>
    </div>
  </section>
</template>
<script>
export default {
  name: 'CallToActionComponent',
}
</script>
