<template>
  <div v-bind:id="id" class="modal">
    <form v-bind:id="`${id}-form`" v-on:submit.prevent="add_parser">
      <div class="modal-content">
        <h4>Add parser</h4>
        <div class="row">
          <div class="input-field col s12">
            <input v-model="name" id="name" type="text" class="validate" required>
            <label for="name">Name</label>
          </div>
        </div>
        <div class="file-field input-field">
          <div class="btn">
            <span>File</span>
            <input type="file" ref="model" v-on:change="on_file_change" required>
          </div>
          <div class="file-path-wrapper">
            <input class="file-path validate" type="text">
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
export default {
  name: 'add_parser',
  props: [ 'id' ],
  data() {
    return {
      name: '',
      parser: null,
    };
  },
  methods: {
    on_file_change: function() {
      if (this.$refs.model.files.length > 0) {
        this.model = this.$refs.model.files[0];
      }
    },
    add_parser: function() {
      this.$store.dispatch('add_parser', {
        name: this.name,
        parser: this.parser,
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
}
</script>
