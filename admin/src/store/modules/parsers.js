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
