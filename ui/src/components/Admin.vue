<template>
  <div>
    <ProgressDialog :show="inProgress"></ProgressDialog>

  <v-row
    justify="center"
    align="start"
  >
    <v-col cols="12"
           sm="8"
           md="4">

      <v-tabs
        fixed-tabs
        background-color="indigo"
        dark
        icons-and-text
        class="elevation-2"
        centered
        vertical
      >
        <v-tab key="freeze" href="#tab-freeze" @click="getFreeze">
          Freeze
          <v-icon>mdi-lock</v-icon>
        </v-tab>
        <v-tab key="unfreeze" href="#tab-unfreeze" @click="getUnFreeze">
          UnFreeze
          <v-icon>mdi-lock-open-variant</v-icon>
        </v-tab>
        <v-tab key="clawback" href="#tab-clawback" @click="getClawback">
          Clawback
          <v-icon>mdi-gauge-empty</v-icon>
        </v-tab>
        <v-tab key="changeAdmin" href="#tab-changeAdmin">
          Admin Key
          <v-icon>mdi-account-key</v-icon>
        </v-tab>
        <v-tab-item
          key="freeze"
          value="tab-freeze"
        >
            <v-form>
              <v-container
                fill-height
                fluid>
                <v-row dense v-for="user in users" :key="user.text">
                  <v-col>
                    <v-checkbox v-model="list" :label="user.text" :value="user.text"></v-checkbox>
                  </v-col>
                </v-row>
              </v-container>
            </v-form>

            <v-card-actions
            >
              <v-col class="text-right">
                <v-btn
                  right
                  v-if="listValid"
                  color="deep-purple lighten-2"
                  text
                  @click="postAction"
                >
                  Submit
                </v-btn>
              </v-col>
            </v-card-actions>
        </v-tab-item>
        <v-tab-item
          key="unfreeze"
          value="tab-unfreeze"
        >
          <v-form>
            <v-container
              fill-height
              fluid>
              <v-row dense v-for="user in users" :key="user.text">
                <v-col>
                  <v-checkbox v-model="list" :label="user.text" :value="user.text"></v-checkbox>
                </v-col>
              </v-row>
            </v-container>
          </v-form>

          <v-card-actions
          >
            <v-col class="text-right">
              <v-btn
                right
                v-if="listValid"
                color="deep-purple lighten-2"
                text
                @click="postAction"
              >
                Submit
              </v-btn>
            </v-col>
          </v-card-actions>
        </v-tab-item>
        <v-tab-item
          key="clawback"
          value="tab-clawback"
        >
          <v-form>
            <v-container
              fill-height
              fluid>
              <v-row dense v-for="user in users" :key="user.text">
                <v-col>
                  <v-checkbox v-model="list" :label="user.text" :value="user.text"></v-checkbox>
                </v-col>
              </v-row>
            </v-container>
          </v-form>

          <v-card-actions
          >
            <v-col class="text-right">
              <v-btn
                right
                v-if="listValid"
                color="deep-purple lighten-2"
                text
                @click="postAction"
              >
                Submit
              </v-btn>
            </v-col>
          </v-card-actions>
        </v-tab-item>
        <v-tab-item
          key="changeAdmin"
          value="tab-changeAdmin"
        >
          <v-form v-model="keyValid">
            <v-container
              fill-height
              fluid>
              <v-row>
                <v-col>
                  <v-text-field
                    v-model="key"
                    :rules="keyRules"
                    label="New Private Key"
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
                Submit
              </v-btn>
            </v-col>
          </v-card-actions>
        </v-tab-item>
      </v-tabs>
    </v-col>
  </v-row>
  </div>
</template>

<script>
  import Cookie from 'js-cookie'
  import ProgressDialog from './ProgressDialog'
  import router from '../router'
  import axios from 'axios'
  import {bus} from '../main'
  import * as Proto from '../proto/messages_pb'
  import Utils from '../utils'

  const {
    Ed25519PrivateKey
  } = require('@hashgraph/sdk')

  export default {
    name: 'Admin',
    components: {
      ProgressDialog
    },
    data () {
      return {
        getRestAPI: 'http://' + window.location.hostname + ':' + process.env.GET_PORT,
        postRestAPI: 'http://' + window.location.hostname + ':' + process.env.POST_PORT,
        inProgress: false,
        users: [],
        list: [],
        getFrozen: false,
        activity: 'freeze',
        keyValid: false,
        key: '',
        keyRules: [
          v => v.toString().length > 0 || 'Key is required'
        ]
      }
    },
    methods: {
      home: function () {
        router.replace('/')
      },
      getFreeze: function () {
        this.getFrozen = false
        this.getUsers()
        this.activity = 'freeze'
      },
      getUnFreeze: function () {
        this.getFrozen = true
        this.getUsers()
        this.activity = 'unfreeze'
      },
      getClawback: function () {
        this.getFrozen = true
        this.getUsers()
        this.activity = 'clawback'
      },
      getUsers: function () {
        this.list = []
        let url = ''
        if (this.getFrozen) {
          url = this.getRestAPI.concat('/v1/token/users/frozen')
        } else {
          url = this.getRestAPI.concat('/v1/token/users/unfrozen')
        }

        axios.get(url)
          .then(response => {
            this.users = response.data
          })
          .catch(e => {
            console.log(e)
          })
      },
      postAction: function () {
        const privateKey = Cookie.get('userKey')
        this.list.forEach(username => {
          const body = {}
          let primitive = new Proto.Primitive()
          if (this.activity === 'freeze') {
            const freezeProto = new Proto.Freeze()
            freezeProto.setAccount(username)
            primitive.setHeader(Utils.primitiveHeader(freezeProto.serializeBinary(), privateKey))
            primitive.setFreeze(freezeProto)
          } else if (this.activity === 'unfreeze') {
            const unFreezeProto = new Proto.UnFreeze()
            unFreezeProto.setAccount(username)
            primitive.setHeader(Utils.primitiveHeader(unFreezeProto.serializeBinary(), privateKey))
            primitive.setUnfreeze(unFreezeProto)
          } else {
            const clawbackProto = new Proto.Clawback()
            clawbackProto.setAccount(username)
            primitive.setHeader(Utils.primitiveHeader(clawbackProto.serializeBinary(), privateKey))
            primitive.setClawback(clawbackProto)
          }

          const serializedPrimitive = primitive.serializeBinary()

          const base64 = btoa(String.fromCharCode.apply(null, serializedPrimitive))
          body.primitive = base64

          this.inProgress = true

          axios.post(this.postRestAPI.concat('/v1/token/transaction'), body)
            .then(response => {
              this.inProgress = false
              console.log(response.data.message)
              if (response.data.status) {
                bus.$emit('showSuccess', response.data.message)
              } else {
                bus.$emit('showError', response.data.message)
              }
            })
            .catch(e => {
              this.inProgress = false
              console.log(e)
              bus.$emit('showError', e)
            })
        })
      },
      postKey: function () {
        const privateKey = Cookie.get('userKey')

        const body = {}
        let primitive = new Proto.Primitive()
        const adminKeyProto = new Proto.AdminKeyUpdate()
        // get public key from private key
        const privKey = Ed25519PrivateKey.fromString(this.key)
        const pubKey = privKey.publicKey.toString()

        adminKeyProto.setNewpublickey(pubKey)
        primitive.setHeader(Utils.primitiveHeader(adminKeyProto.serializeBinary(), privateKey))
        primitive.setAdminkeyupdate(adminKeyProto)

        const serializedPrimitive = primitive.serializeBinary()

        const base64 = btoa(String.fromCharCode.apply(null, serializedPrimitive))
        body.primitive = base64

        this.inProgress = true

        axios.post(this.postRestAPI.concat('/v1/token/transaction'), body)
          .then(response => {
            this.inProgress = false
            console.log(response.data.message)
            if (response.data.status) {
              bus.$emit('showSuccess', response.data.message)
              Cookie.set('userKey', this.key)
            } else {
              bus.$emit('showError', response.data.message)
            }
          })
          .catch(e => {
            this.inProgress = false
            console.log(e)
            bus.$emit('showError', e)
          })
      }
    },
    created () {
      if (Cookie.get('isAdmin') === '') {
        router.go(0)
      }

      this.getUsers()

      bus.$on('refresh', (message) => {
        this.getUsers()
      })
    },
    computed: {
      listValid () {
        if (this.list.length > 0) {
          return 'true'
        } else {
          return ''
        }
      }
    }
  }
</script>

<style scoped>
</style>
