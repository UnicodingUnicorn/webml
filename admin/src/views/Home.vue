<template>
  <div class="container">
    <div class="row">
      <div class="col s12">
        <div class="card">
          <div class="card-content">
            <span class="card-title" v-if="currentModel">{{currentModel.name}}</span>
            <span class="card-title" v-else>No model selected</span>
          </div>
          <div class="card-action" v-if="currentModel">
            <a class="modal-trigger" v-bind:href="`#add-dataset-${currentModel.id}`">Add dataset</a>
          </div>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col s6">
        <div class="card">
          <div class="card-content">
            <span class="card-title">Models</span>
            <table v-if="models.length > 0">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>No. datasets</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="model in models" v-bind:key="model.id">
                  <td v-on:click="changeModel(model.id)" style="cursor:pointer">{{model.name}}</td>
                  <td v-on:click="changeModel(model.id)" style="cursor:pointer">{{model.data.length}}</td>
                </tr>
              </tbody>
            </table>
            <i v-else>No models found</i>
          </div>
        </div>
      </div>
      <div class="col s6">
        <div class="card">
          <div class="card-content">
            <span class="card-title">Parsers</span>
            <table v-if="parsers.length > 0">
              <thead>
                <tr>
                  <th>Name</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="parser in parsers" v-bind:key="parser.id">
                  <td>{{parser.name}}</td>
                </tr>
              </tbody>
            </table>
            <i v-else>No parsers found</i>
          </div>
        </div>
      </div>
    </div>
    <div class="fixed-action-btn">
      <a class="btn-floating btn-large red">
        <i class="large material-icons">add</i>
      </a>
      <ul>
        <li><a class="btn-floating green modal-trigger" title="Add model" href="#add-model"><i class="material-icons">layers</i></a></li>
        <li><a class="btn-floating blue modal-trigger" title="Add parser" href="#add-parser"><i class="material-icons">tune</i></a></li>
      </ul>
    </div>
    <AddParser id="add-parser" />
    <AddModel id="add-model" />
    <AddDataset v-for="model in models" v-bind:key="model.id" v-bind:id="`add-dataset-${model.id}`" v-bind:model="model.id" />
  </div>
</template>

<script>
import AddModel from '@/components/AddModel.vue'
import AddParser from '@/components/AddParser.vue'
import AddDataset from '@/components/AddDataset.vue'

import { mapGetters } from 'vuex'

export default {
  name: 'home',
  data() {
    return {
      currentID: null,
    };
  },
  components: {
    AddModel,
    AddParser,
    AddDataset,
  },
  computed: {
    ...mapGetters([
      'models',
      'parsers',
    ]),
    currentModel: function() {
      const model = this.models.find(model => model.id == this.currentID);
      return model != undefined ? model : null;
    },
  },
  methods: {
    changeModel(id) {
      this.currentID = id;
    },
  },
  mounted: function() {
    this.$store.dispatch('init').then(() => {
      console.log('stuff loaded');
      console.log(this.parsers);

      // Load modals
      const modalElems = document.querySelectorAll('.modal');
      // eslint-disable-next-line
      M.Modal.init(modalElems);
    }).catch(err => {
      let msg = '';
      if (err == 400) {
        msg = 'An error has occurred';
      } else if (err == 404) {
        msg = 'A resource cannot be found';
      } else if (err == 500) {
        msg = 'A server error has occurred';
      }
      //eslint-disable-next-line
      M.toast({ html: msg });
    });
  },
}
</script>
