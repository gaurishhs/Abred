<template>
  <div>
    <div class="context">
      <h1 class="text-black py-8">Abred</h1>
      <p class="text-black text-center">
        The only bot you will need for managing temporary voice channels.
      </p>
      <br />
      <div class="flex justify-center">
        <a href="/invite" class=""
          ><button
            class="inline-flex text-white bg-gray-800 border-0 py-2 px-6 focus:outline-none hover:bg-red-600 rounded text-lg"
          >
            Invite
          </button></a
        >
        <a href="/manage" class=""
          ><button
            class="ml-3 inline-flex text-gray-400 bg-gray-800 border-0 py-2 px-6 focus:outline-none hover:bg-green-600 hover:text-white rounded text-lg"
          >
            Manage
          </button></a
        >
        <a href="/support" class=""
          ><button
            class="ml-4 inline-flex text-gray-400 bg-gray-800 border-0 py-2 px-6 focus:outline-none hover:bg-pink-600 hover:text-white rounded text-lg"
          >
            Support
          </button></a
        >
      </div>
    </div>
    <div class="area">
      <ul class="circles">
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
      </ul>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import '~/styles/home.scss'
export default Vue.extend({
  name: 'Index',
  data() {
    let open = false

    return { open }
  },
  loading: {
    color: 'blue',
  },
  head: {
    title: 'Home | Abred',
    meta: [
      {
        charset: 'utf-8',
      },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      {
        hid: 'description',
        name: 'description',
        content:
          'Abred the only discord bot you need to manage temporary voice channels in your guild.',
      },
    ],
  },
z
})
</script>

<style lang="scss">
.slide-left {
  position: relative;
  overflow: hidden;
  width: 100px;
  height: 100px;
}

.dash-btn {
  background-color: rgb(16 185 129);
  &:hover {
    background-color: rgb(6 95 70);
  }
}

#slide-left {
  position: absolute;
  left: -100px;
  width: 100px;
  height: 100px;
  align-content: center;
  -webkit-animation: slide 0.5s forwards;
  -webkit-animation-delay: 2s;
  animation: slide 0.5s forwards;
  animation-delay: 2s;
}

@-webkit-keyframes slide {
  100% {
    left: 0;
  }
}

@keyframes slide {
  100% {
    left: 0;
  }
}
</style>
