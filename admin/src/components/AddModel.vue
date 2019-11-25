<template>
  <div v-bind:id="id" class="modal">
    <form v-bind:id="`${id}-form`" v-on:submit.prevent="add_model">
      <div class="modal-content">
        <h4>Add model</h4>
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
  name: 'add_model',
  props: [ 'id' ],
  data() {
    return {
      name: '',
      model: null,
    };
  },
  methods: {
    on_file_change: function() {
      if (this.$refs.model.files.length > 0) {
        this.model = this.$refs.model.files[0];
      }
    },
    add_model: function() {
      this.$store.dispatch('add_model', {
        name: this.name,
        model: this.model,
      }).then(() => {
        document.getElementById(`${this.id}-form`).reset();
        // eslint-disable-next-line
        M.Modal.getInstance(document.getElementById(this.id)).close();
      }).catch(e => {
        let msg = '';
        if (e == 400) {
          msg = 'incorrect data supplied';
        } else if (e == 409) {
          msg = 'id conflict';
        } else if (e == 500) {
          msg = 'server error encountered';
        }
        // eslint-disable-next-line
        M.toast({ html: msg });
      });
    },
    reset_form: function() {
      document.getElementById(`${this.id}-form`).reset();
    },
  },
}
</script>
