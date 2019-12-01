<template>
  <div v-bind:id="id" class="modal">
    <form v-bind:id="`${id}-form`" v-on:submit.prevent="add_dataset">
      <div class="modal-content">
        <h4>Add Dataset</h4>
        <div class="row">
          <div class="input-field col s12">
            <input v-model="shape" id="name" type="text" class="validate" required>
            <label for="shape">Shape (JSON float array) </label>
          </div>
        </div>
        <div class="row">
          <div class="file-field input-field col s12">
            <div class="btn">
              <span>Data File</span>
              <input type="file" ref="data" v-on:change="on_data_file_change" required>
            </div>
            <div class="file-path-wrapper">
              <input class="file-path validate" type="text">
            </div>
          </div>
          <div class="input-field col s12">
            <select required v-model="data_parser">
              <option value="" disabled selected>Choose data parser</option>
              <option v-for="parser in parsers" v-bind:value="parser.id" v-bind:key="parser.name">{{parser.name}}</option>
            </select>
            <label>Data parser</label>
          </div>
        </div>
        <div class="row">
          <div class="file-field input-field col s12">
            <div class="btn">
              <span>Labels File</span>
              <input type="file" ref="labels" v-on:change="on_labels_file_change" required>
            </div>
            <div class="file-path-wrapper">
              <input class="file-path validate" type="text">
            </div>
          </div>
          <div class="input-field col s12">
            <select required v-model="labels_parser">
              <option value="" disabled selected>Choose labels parser</option>
              <option v-for="parser in parsers" v-bind:value="parser.id" v-bind:key="parser.name">{{parser.name}}</option>
            </select>
            <label>Labels parser</label>
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="waves-effect waves-green btn-flat">Agree</button>
        <a href="#" v-on:click="reset_form" class="modal-close waves-effect waves-red btn-flat">Close</a>
      </div>
    </form>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'

export default {
  name: 'add_dataset',
  props: [ 'id', 'model' ],
  data() {
    return {
      name: '',
      shape: '',
      data: null,
      data_parser: '',
      labels: null,
      labels_parser: '',
    };
  },
  computed: {
    ...mapGetters([
      'parsers',
    ]),
  },
  methods: {
    on_data_file_change: function() {
      if (this.$refs.data.files.length > 0) {
        this.data = this.$refs.data.files[0];
      }
    },
    on_labels_file_change: function() {
      if (this.$refs.labels.files.length > 0) {
        this.labels = this.$refs.labels.files[0];
      }
    },
    add_dataset: function() {
      this.$store.dispatch('add_dataset', {
        model: this.model,
        data: this.data,
        labels: this.labels,
        shape: this.shape,
        dataparser: this.data_parser,
        labelparser: this.labels_parser,
      }).then(() => {
        document.getElementById(`${this.id}-form`).reset();
        // eslint-disable-next-line
        M.Modal.getInstance(document.getElementById(this.id)).close();
      }).catch(e => {
        // eslint-disable-next-line
        M.toast({ html: e });
      });
    },
    reset_form: function() {
      document.getElementById(`${this.id}-form`).reset();
    },
  },
  updated: function() {
    const elems = document.querySelectorAll('select');
    // eslint-disable-next-line
    M.FormSelect.init(elems);
  },
}
</script>
