import Vue from 'vue'
import Vuex from 'vuex'

import models from './modules/models';
import parsers from './modules/parsers';

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    models,
    parsers,
  },
  actions: {
    init({ dispatch }) {
      return dispatch('init_models').then(() => {
        return dispatch('init_parsers');
      }).catch(e => {
        if (e == 400) {
          throw 'An error has occurred';
        } else if (e == 404) {
          throw 'A resource cannot be found';
        } else if (e == 500) {
          throw 'A server error has occurred';
        } else if (typeof e === 'string') {
          throw e;
        } else {
          throw 'A client error has occurred';
        }
      });
    },
  },
});
