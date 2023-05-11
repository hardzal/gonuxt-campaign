<template>
  <div class="container mx-auto h-screen flex justify-center items-center">
    <div class="w-full lg:w-1/3 px-10 lg:px-0">
      <div class="flex justify-center items-center mx-auto mt-6 mb-8">
        <img src="/success-illustration.svg" alt="" class="w-full" />
      </div>
      <h2 class="font-medium mb-3 text-3xl text-center">
        Oops! Something went wrong.
      </h2>
      <p v-if="error.StatusCode == 404" class="text-center font-light">
        The page that you requested doesn't exist at the moment
      </p>
      <p v-else class="text-center font-light">An error occured</p>
      <div class="mb-4 mt-6">
        <div class="mb-3">
          <button
            class="block w-full bg-orange-button hover:bg-green-button text-white font-semibold px-6 py-4 text-lg rounded-full"
            @click="$router.push({ path: '/' })"
          >
            Start Explore
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'success',
  //   props: ['error'],
}
</script>
