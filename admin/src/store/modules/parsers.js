import * as types from '../mutation-types';

import parserApi from '../../resource/parser';

export default {
  state: {
    parsers: [],
  },
  getters: {
    parsers: state => state.parsers,
    parser: state => id => state.parsers.find(p => p.id === id) || null,
  },
  actions: {
    init_parsers({ commit }) {
      return parserApi.get_parsers().then(parserids => {
        return Promise.all(parserids.map(parserid => {
          return parserApi.get_parser(parserid);
        })).then(parsers => {
          return parsers;
        });
      }).then(parsers => {
        commit(types.INIT_PARSERS, parsers);
      });
    },
    add_parser({ commit }, { name, parser }) {
      return parserApi.put_parser(parser, name).then(() => {
        commit(types.ADD_PARSER, {
          name,
          parser,
        });
      }).catch(e => {
        if (e == 400) {
          throw 'Incorrect data supplied';
        } else if (e == 409) {
          throw 'ID conflict';
        } else if (e == 500) {
          throw 'Server error encountered';
        } else if (typeof e === 'string') {
          throw e;
        } else {
          throw 'A client error has occurred';
        }
      });
    },
  },
  mutations: {
    [types.INIT_PARSERS](state, parsers) {
      state.parsers = parsers;
    },
    [types.ADD_PARSER](state, parser) {
      state.parsers.push(parser);
    },
  },
}
