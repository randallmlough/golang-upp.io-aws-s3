// Rollup plugins
import babel from 'rollup-plugin-babel';
import resolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';
import replace from 'rollup-plugin-replace';
import { terser } from "rollup-plugin-terser";

export default {
    input: 'src/main.js',
    output: {
        file: 'assets/js/bundle.js',
        format: 'iife',
        name: 'app'
    },
    plugins: [
        resolve({
            mainFields: [
                'main',
                'jsnext:main',
                'browser'
            ]
        }),
        commonjs(),
        babel({
            exclude: 'node_modules/**',
        }),
        replace({
            ENV: JSON.stringify(process.env.NODE_ENV || 'development'),
        }),
        (process.env.NODE_ENV === 'production' && terser()),
    ],
};