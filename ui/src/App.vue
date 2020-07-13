<template>
  <v-app id="app">
    <router-view></router-view>
    <v-snackbar
      v-model="snackbar"
      :bottom=true
      :color="snackAlert"
      :multi-line=false
      :timeout=5000
      :vertical=true
    >
      {{ snackMessage }}
      <v-btn
        dark
        text
        @click="snackbar = false"
      >
        Close
      </v-btn>
    </v-snackbar>
  </v-app>
</template>

<script>
import Cookie from 'js-cookie'
import { bus } from './main'

export default {
  name: 'app',
  data: function () {
    return {
      snackbar: false,
      snackMessage: '',
      snackAlert: 'alert'
    }
  },
  methods: {
  },
  created: function () {
    this.$socketClient.onOpen = () => {
      if ((typeof (Cookie.get('userName')) !== 'undefined') && (Cookie.get('userName') !== '')) {
        // register web socket with username
        console.log('Establishing Socket Connection')
        this.$socketClient.sendObj({userId: Cookie.get('userName')})
      }
    }
    this.$socketClient.onMessage = (msg) => {
      this.snackbar = false
      const notification = JSON.parse(msg.data)
      if (notification.error) {
        this.snackAlert = 'error'
      } else {
        this.snackAlert = 'success'
      }
      this.snackMessage = notification.message
      this.snackbar = true
      bus.$emit('refresh')
    }
    this.$socketClient.onClose = (msg) => {
      console.log('socket closed')
    }
    this.$socketClient.onError = (msg) => {
      console.log('socket error')
    }
  },
  mounted () {
    bus.$on('showSuccess', (message) => {
      this.snackAlert = 'success'
      this.snackMessage = message
      this.snackbar = true
    })
    bus.$on('showError', (message) => {
      this.snackAlert = 'error'
      this.snackMessage = message
      this.snackbar = true
    })
    bus.$on('hideSnackbar', () => {
      this.snackbar = false
      this.snackMessage = ''
    })
  }
}
</script>

<style>
body {
  margin: 0;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

main {
  text-align: center;
  margin-top: 40px;
}

header {
  margin: 0;
  height: 56px;
  padding: 0 16px 0 24px;
  background-color: #35495E;
  color: #ffffff;
}

header span {
  display: block;
  position: relative;
  font-size: 20px;
  line-height: 1;
  letter-spacing: .02em;
  font-weight: 400;
  box-sizing: border-box;
  padding-top: 16px;
}
</style>
