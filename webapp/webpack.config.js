const path = require('path');

module.exports = {
    entry: './src/index.tsx',
    mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'main.js',
        library: {
            type: 'window',
        },
    },
    externals: {
        react: 'window.React',
        'react-dom': 'window.ReactDOM',
        redux: 'window.Redux',
        'react-redux': 'window.ReactRedux',
        'prop-types': 'window.PropTypes',
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: [/node_modules/, /\.test\.tsx?$/, /test\//],
            },
            {
                test: /\.css$/,
                use: [
                    'style-loader',
                    'css-loader',
                ],
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js', '.jsx', '.css'],
    },
    devtool: process.env.NODE_ENV === 'production' ? false : 'source-map',
};
