<template>
  <div>
    <navbar />
    <div
      wire:loading
      v-if="loading"
      class="fixed top-0 left-0 right-0 bottom-0 w-full h-screen z-50 overflow-hidden bg-gray-700 opacity-75 flex flex-col items-center justify-center"
    >
      <div
        class="loader ease-linear rounded-full border-4 border-t-4 border-gray-200 h-12 w-12 mb-4"
      ></div>
      <h2 class="text-center text-white text-xl font-semibold">Loading...</h2>
      <p class="w-1/3 text-center text-white">
        Hang on tight, Trying to establish a secure connection with server, This
        might take a while, Please do not close this tab.
      </p>
    </div>
    <div v-else>
      <p
        class="text-black text-bold text-center px-5 py-7"
        v-if="this.noconnect"
      >
        You aren't connected to any voice channel currently!
      </p>

      <div v-if="!this.noconnect">
        <p class="text-bold text-center text-base text-black px-5 py-7">
          You are currently connected to {{ this.channelName }} in
          {{ this.guildName }}
        </p>
        <br />
        <div>
          <h1 class="text-pink-500 text-center text-base text-2xl">Info:</h1>
          <br />
          <p class="text-center">Current State: {{ this.getState() }}</p>
          <p class="text-center">User Limit: {{ this.userlimit }}</p>
          <p v-if="this.pushtotalk" class="text-center">
            Push To Talk: <span class="material-icons">done</span>
          </p>
          <p v-else class="text-center">
            Push To Talk: <span class="material-icons">close</span>
          </p>
        </div>
        <br />
        <div class="justify-center items-center content-center">
          <h1 class="text-red-500 text-bold text-2xl text-center">Actions:</h1>
          <center>
            <button
              type="button"
              v-on:click="this.increaseLimit"
              class="text-white action-btn bg-gradient-to-br from-green-400 to-blue-600 hover:bg-gradient-to-bl focus:ring-4 focus:outline-none focus:ring-green-200 dark:focus:ring-green-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center mr-2 mb-2"
              :disabled="this.incdisabled"
            >
              Increase Limit
            </button>
            <button
              type="button"
              v-on:click="this.decreaseLimit"
              class="text-white action-btn bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center mr-2 mb-2"
              :disabled="this.decdisabled"
            >
              Decrease Limit
            </button>
            <button
              type="button"
              v-on:click="this.togglePushToTalk"
              class="text-white action-btn bg-gradient-to-br from-pink-500 to-orange-400 hover:bg-gradient-to-bl focus:ring-4 focus:outline-none focus:ring-pink-200 dark:focus:ring-pink-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center mr-2 mb-2"
              :disabled="this.pttdisabled"
            >
              Toggle Push To Talk
            </button>
          </center>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Navbar from '~/components/Navbar.vue'
import '~/styles/manage.scss'
export default Vue.extend({
  components: { Navbar },
  name: 'Manage',
  data() {
    let loading: boolean = true
    let channelName: string = ''
    let guildName: string = ''
    let guildId: string = ''
    let channelId: string = ''
    let noconnect: boolean = false
    let locked: boolean = false
    let userlimit: number = 0
    let pushtotalk: boolean = false
    let incdisabled: boolean = false
    let decdisabled: boolean = false
    let pttdisabled: boolean = false
    return {
      loading,
      noconnect,
      pttdisabled,
      locked,
      channelName,
      incdisabled,
      guildName,
      guildId,
      channelId,
      decdisabled,
      userlimit,
      pushtotalk,
    }
  },
  async mounted() {
    if(!this.$auth.user?.id) {
      console.log(this.$auth.user)
      this.$router.push('/login')
    }
    let socket = new WebSocket(
      `wss://api.abred.bar/manage?uid=${this.$auth.user?.id}`
    )

    socket.onopen = () => {
      this.$toasted.success(`Successfully Connected To Server`, {
        duration: 5000,
        position: 'bottom-right',
      })
    }

    socket.onmessage = (data) => {
      this.loading = false
      if (data.data == 'none') {
        this.noconnect = true
        return
      }
      this.noconnect = false
      const vcData: any = JSON.parse(data.data)
      this.pushtotalk = vcData.pushToTalk
      this.locked = vcData.locked
      this.userlimit = vcData.userLimit
      this.channelName = vcData.channelName
      this.guildName = vcData.guildName
      this.guildId = vcData.guildId
      this.channelId = vcData.channelId
    }

    socket.onclose = () => {
      this.$toasted.error('Disconnected From The Server', {
        duration: 5000,
        position: 'bottom-right',
      })
      this.loading = true
    }
  },
  head: {
    title: 'Manage | Abred',
  },
  methods: {
    getState() {
      if (this.locked == true) {
        return 'locked'
      } else {
        return 'unlocked'
      }
    },
    increaseLimit() {
      this.incdisabled = true
      this.$axios({
        url: 'https://api.abred.bar/inc-limit',
        method: 'POST',
        params: {
          gid: this.guildId,
          uid: this.$auth.user?.id,
          cid: this.channelId,
        },
      })
        .then((response) => {
          if (response.data.success) {
            this.$toasted.success(
              `Successfully Increased User Limit of ${this.channelName}`,
              {
                position: 'bottom-right',
                duration: 5000,
              }
            )
          }
          setTimeout(() => {
            this.incdisabled = false
          }, 10000)
        })
        .catch(() => {
          this.incdisabled = false
          this.$toasted.error(
            `Could not increase user limit of ${this.channelName}`,
            {
              duration: 5000,
              position: 'bottom-right',
            }
          )
        })
    },
    decreaseLimit() {
      this.decdisabled = true
      this.$axios({
        url: 'https://api.abred.bar/dec-limit',
        method: 'POST',
        params: {
          gid: this.guildId,
          uid: this.$auth.user?.id,
          cid: this.channelId,
        },
      })
        .then((response) => {
          if (response.data.success) {
            this.$toasted.success(
              `Successfully Decreased User Limit of ${this.channelName}`,
              {
                position: 'bottom-right',
                duration: 5000,
              }
            )
          }
          setTimeout(() => {
            this.decdisabled = false
          }, 10000)
        })
        .catch(() => {
          this.decdisabled = false
          this.$toasted.error(
            `Could not decrease user limit of ${this.channelName}`,
            {
              duration: 5000,
              position: 'bottom-right',
            }
          )
        })
    },
    togglePushToTalk() {
      this.pttdisabled = true
      this.$axios({
        method: 'POST',
        url: 'https://api.abred.bar/ptt',
        params: {
          gid: this.guildId,
          uid: this.$auth.user?.id,
          cid: this.channelId,
        },
      })
        .then((response) => {
          if (response.data.success) {
            this.$toasted.success(
              `Successfully toggled push to talk for ${this.channelName}`,
              {
                position: 'bottom-right',
                duration: 5000,
              }
            )
            setTimeout(() => {
              this.pttdisabled = false
            }, 10000)
          }
        })
        .catch(() => {
          this.$toasted.error(
            `Could not toggle push to talk for ${this.channelName}`,
            {
              duration: 5000,
              position: 'bottom-right',
            }
          )
          this.pttdisabled = false
        })
    },
  },
})
</script>
