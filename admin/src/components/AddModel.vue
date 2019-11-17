<template>
  <div v-bind:id="id" class="modal">
    <form>
      <div class="modal-content">
        <h4>Add model</h4>
        <div class="row">
          <div class="input-field col s12">
            <input id="name" type="text" class="validate" required>
            <label for="name">Name</label>
          </div>
        </div>
        <div class="file-field input-field">
          <div class="btn">
            <span>File</span>
            <input type="file">
          </div>
          <div class="file-path-wrapper">
            <input class="file-path validate" type="text">
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <a href="#!" class="modal-close waves-effect waves-green btn-flat">Agree</a>
      </div>
    </form>
  </div>
</template>
<script>
export default {
  name: 'add_model',
  props: [ 'id' ],
}
</script>
