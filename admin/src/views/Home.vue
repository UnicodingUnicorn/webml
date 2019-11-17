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
            <a href="#">Add dataset</a>
          </div>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col s6">
        <div class="card-panel">
          <table v-if="models.length > 0">
            <thead>
              <tr>
                <th>Name</th>
                <th>No. datasets</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="model in models" v-bind:key="model.id">
                <td>{{model.name}}</td>
                <td>{{model.data.length}}</td>
              </tr>
            </tbody>
          </table>
          <i v-else>No models found</i>
        </div>
      </div>
      <div class="col s6">
        <div class="card-panel">
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
  </div>
</template>

<script>
import AddModel from '@/components/AddModel.vue'
import AddParser from '@/components/AddParser.vue'
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
  },
  computed: {
    ...mapGetters([
      'models',
      'parsers',
    ]),
    currentModel: function() {
      const model =  this.models.find(model => model.id == this.currentID);
      return model != undefined ? model : null;
    },
  },
  mounted: function() {
    this.$store.dispatch('init').then(() => {
      console.log('stuff loaded');
      console.log(this.parsers);
    });
  },
}
</script>
