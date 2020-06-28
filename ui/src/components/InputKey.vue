<template>
  <div>
    <ProgressDialog :show="inProgress"></ProgressDialog>

    <v-toolbar dense>
      <v-btn icon @click="home">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>Private Key Input</v-toolbar-title>
    </v-toolbar>
    <v-sheet
      height="100%"
    >

  <v-row
    justify="center"
    align="start"
  >
    <v-col cols="12"
           sm="8"
           md="4">

      <v-card
        class="mx-auto"
        max-width="344"
        raised
        title="Private Key Input"
      >
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title class="headline">Private Key Input</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-form v-model="keyValid">
          <v-container
            fill-height
            fluid>
            <v-row>
              <v-col>
                <v-text-field
                  v-model="key"
                  :rules="keyRules"
                  label="Private Key"
                  type="password"
                  required
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-form>

        <v-card-actions
        >
          <v-col class="text-right">
            <v-btn
              right
              v-if="keyValid"
              color="deep-purple lighten-2"
              text
              @click="postKey"
            >
              Set Key
            </v-btn>
          </v-col>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
    </v-sheet>
  </div>
</template>

<script>
  import Cookie from 'js-cookie'
  import ProgressDialog from './ProgressDialog'
  import router from '../router'

  export default {
    name: 'InputKey',
    components: {
      ProgressDialog
    },
    data () {
      return {
        restAPI: 'http://' + window.location.hostname + ':' + process.env.HOST_PORT,
        inProgress: false,
        keyValid: false,
        key: '',
        keyRules: [
          v => v.toString().length > 0 || 'Key is required'
        ]
      }
    },
    created () {
    },
    methods: {
      home: function () {
        router.replace('/')
      },
      postKey: function () {
        Cookie.set('userName', 'Admin')
        Cookie.set('userKey', this.key)
        Cookie.set('friendlyName', 'Admin')
        router.replace('/')
      }
    }
  }
</script>

<style scoped>
</style>
