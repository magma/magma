/* eslint-disable */
// Generated automatically by nearley, version 2.19.0
// http://github.com/Hardmath123/nearley
(function () {
function id(x) { return x[0]; }

const {lexer} = require('../PromQLTokenizer');
const {AggregationOperation, BinaryOperation, Function, InstantSelector, Label, Labels, RangeSelector, Scalar} = require('../PromQL');
var grammar = {
    Lexer: lexer,
    ParserRules: [
    {"name": "EXPRESSION", "symbols": ["METRIC_SELECTOR"], "postprocess": id},
    {"name": "EXPRESSION", "symbols": ["FUNCTION"], "postprocess": id},
    {"name": "EXPRESSION", "symbols": ["AGGREGATION"], "postprocess": id},
    {"name": "EXPRESSION", "symbols": ["BINARY_OPERATION"], "postprocess": id},
    {"name": "EXPRESSION", "symbols": ["SCALAR"], "postprocess": id},
    {"name": "METRIC_SELECTOR", "symbols": ["INSTANT_SELECTOR"], "postprocess": id},
    {"name": "METRIC_SELECTOR", "symbols": ["RANGE_SELECTOR"], "postprocess": id},
    {"name": "INSTANT_SELECTOR", "symbols": ["IDENTIFIER", "label_selector"], "postprocess": ([id, labels]) => new InstantSelector(id, labels)},
    {"name": "INSTANT_SELECTOR", "symbols": ["IDENTIFIER"], "postprocess": ([id]) => new InstantSelector(id)},
    {"name": "INSTANT_SELECTOR", "symbols": ["label_selector"], "postprocess": ([labels]) => new InstantSelector("", labels)},
    {"name": "IDENTIFIER", "symbols": [(lexer.has("word") ? {type: "word"} : word)], "postprocess": ([id]) => id.text},
    {"name": "RANGE_SELECTOR", "symbols": ["INSTANT_SELECTOR", "duration"], "postprocess": ([selector, duration]) => new RangeSelector(selector, duration)},
    {"name": "BINARY_OPERATION", "symbols": ["EXPRESSION", "bin_op", "EXPRESSION"], "postprocess": ([lh, op, rh]) => new BinaryOperation(lh, rh, op)},
    {"name": "FUNCTION", "symbols": ["func_name", (lexer.has("lParen") ? {type: "lParen"} : lParen), "func_params", (lexer.has("rParen") ? {type: "rParen"} : rParen)], "postprocess": ([funcName, _lParen, params, _rParen]) => new Function(funcName, params)},
    {"name": "AGGREGATION", "symbols": ["agg_op", (lexer.has("lParen") ? {type: "lParen"} : lParen), "func_params", (lexer.has("rParen") ? {type: "rParen"} : rParen)], "postprocess": ([aggOp, _lParen, params, _rParen]) => new AggregationOperation(aggOp, params)},
    {"name": "func_params", "symbols": ["func_params", (lexer.has("comma") ? {type: "comma"} : comma), "parameter"], "postprocess": ([existingParams, _comma, newParam]) => [...existingParams, newParam]},
    {"name": "func_params", "symbols": ["parameter"], "postprocess": d => [d[0]]},
    {"name": "parameter", "symbols": ["SCALAR"], "postprocess": id},
    {"name": "parameter", "symbols": ["string"], "postprocess": id},
    {"name": "parameter", "symbols": ["EXPRESSION"], "postprocess": id},
    {"name": "label_selector", "symbols": [(lexer.has("lBrace") ? {type: "lBrace"} : lBrace), "label_match_list", (lexer.has("rBrace") ? {type: "rBrace"} : rBrace)], "postprocess": ([_lBrace, labels, _rBrace]) => {const ret = new Labels(); labels.forEach(l => ret.addLabel(l.name, l.value, l.operator)); return ret}},
    {"name": "label_selector", "symbols": [(lexer.has("lBrace") ? {type: "lBrace"} : lBrace), (lexer.has("rBrace") ? {type: "rBrace"} : rBrace)], "postprocess": d => new Labels()},
    {"name": "label_match_list", "symbols": ["label_match_list", (lexer.has("comma") ? {type: "comma"} : comma), "label_matcher"], "postprocess": ([existingLabels, _, newLabel]) => [...existingLabels, newLabel]},
    {"name": "label_match_list", "symbols": ["label_matcher"], "postprocess": d => [d[0]]},
    {"name": "label_matcher", "symbols": ["IDENTIFIER", "label_op", "string"], "postprocess": ([name, op, value]) => new Label(name, value, op)},
    {"name": "string", "symbols": [(lexer.has("string") ? {type: "string"} : string)], "postprocess": d => d[0].value},
    {"name": "label_op", "symbols": [(lexer.has("labelOp") ? {type: "labelOp"} : labelOp)], "postprocess": d => d[0].value},
    {"name": "bin_op", "symbols": [(lexer.has("binOp") ? {type: "binOp"} : binOp)], "postprocess": d => d[0].value},
    {"name": "agg_op", "symbols": [(lexer.has("aggOp") ? {type: "aggOp"} : aggOp)], "postprocess": d => d[0].value},
    {"name": "func_name", "symbols": [(lexer.has("functionName") ? {type: "functionName"} : functionName)], "postprocess": d => d[0].value},
    {"name": "duration", "symbols": [(lexer.has("lBracket") ? {type: "lBracket"} : lBracket), "range", (lexer.has("rBracket") ? {type: "rBracket"} : rBracket)], "postprocess": ([_lBracket, range, _rBracket]) => range},
    {"name": "range", "symbols": [(lexer.has("range") ? {type: "range"} : range)], "postprocess": d => d[0].value},
    {"name": "SCALAR", "symbols": [(lexer.has("scalar") ? {type: "scalar"} : scalar)], "postprocess": d => new Scalar(d[0].value)}
]
  , ParserStart: "EXPRESSION"
}
if (typeof module !== 'undefined'&& typeof module.exports !== 'undefined') {
   module.exports = grammar;
} else {
   window.grammar = grammar;
}
})();
