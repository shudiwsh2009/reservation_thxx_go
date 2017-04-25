var path = require('path');
var webpack = require('webpack');
var CleanWebpackPlugin = require('clean-webpack-plugin');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var AssetsWebpackPlugin = require('assets-webpack-plugin');
var ProgressBarWebpackPlugin = require('progress-bar-webpack-plugin');

var argv = process.argv;
var DEV_HOT = false;

// TODO change to parse argv instead
if (process.env.DEV_HOT) {
	DEV_HOT = true;
}

var production = process.env.NODE_ENV === 'production';
var plugins = [
	new webpack.NoEmitOnErrorsPlugin(),
	new ExtractTextPlugin('[name]-[contenthash].css'),
	new AssetsWebpackPlugin({
		fullPath: false,
		prettyPrint: true,
	}),
	new webpack.LoaderOptionsPlugin({
		debug: !production,
	}),
];

if (production) {
	plugins = plugins.concat([
		new CleanWebpackPlugin('public/bundles'),
		new webpack.optimize.OccurrenceOrderPlugin(),
		new webpack.optimize.CommonsChunkPlugin({
			name: 'common',
			children: true,
			minChunks: 2,
		}),
		new webpack.optimize.MinChunkSizePlugin({
			minChunkSize: 51200, // ~50kb
		}),
		new webpack.optimize.UglifyJsPlugin({
			mangle: true,
			compress: {
				warnings: false, // Suppress uglification warnings
			},
		}),
		new webpack.DefinePlugin({
			__SERVER__: !production,
			__DEVELOPMENT__: !production,
			__DEVTOOLS__: !production,
			'process.env': {
				BABEL_ENV: JSON.stringify(process.env.NODE_ENV),
				NODE_ENV: JSON.stringify(process.env.NODE_ENV),
			},
		}),
	]);
} else {
	plugins = plugins.concat([
		new webpack.DefinePlugin({
			'process.env': {
				BABEL_ENV: JSON.stringify(process.env.NODE_ENV),
				NODE_ENV: JSON.stringify(process.env.NODE_ENV),
			},
		}),
	]);
}

var config = {
	devtool: production ? 'cheap-module-source-map' : 'cheap-module-eval-source-map',
	plugins: plugins,
	entry: {
		entry: path.join(__dirname, 'assets/javascripts/EntryApp.jsx'),
		// student: path.join(__dirname, 'assets/javascripts/StudentApp.jsx'),
		// teacher: path.join(__dirname, 'assets/javascripts/TeacherApp.jsx'),
	},
	output: {
		path: path.join(__dirname, 'public/bundles'),
		publicPath: '/assets/bundles/',
		filename: '[name]-[hash].js',
		chunkFilename: '[name]-[chunkhash].js',
	},
	resolve: {
		extensions: ['.js', '.jsx', '.css', '.png'],
		alias: {
			'#pages': path.join(__dirname, 'assets/javascripts/pages'),
			'#forms': path.join(__dirname, 'assets/javascripts/forms'),
			'#coms': path.join(__dirname, 'assets/javascripts/components'),
			'#models': path.join(__dirname, 'assets/javascripts/models'),
			'#utils': path.join(__dirname, 'assets/javascripts/utils'),
			'#imgs': path.join(__dirname, 'assets/images'),
		},
	},
	module: {
		loaders: [
			{ test: /\.css$/, loader: ExtractTextPlugin.extract({ fallback: 'style-loader', use: 'css-loader' }) },
			{ test: /\.scss$/, loader: ExtractTextPlugin.extract({ fallback: 'style-loader', use: 'css-loader!sass-loader' }) },
			{ test: /\.html$/, loader: 'html-loader' },
			{ test: /\.(png|gif|svg)$/, loader: 'url-loader?name=[name]@[hash].[ext]&limit=5000' },
			{ test: /\.(pdf|ico|jpg|eot|otf|woff|ttf|mp4|webm)$/, loader: 'file-loader?name=[name]@[hash].[ext]' },
			{
				test: /\.jsx?$/,
				loader: "babel-loader",
				query: {
					presets: ['es2015', 'react', 'stage-0'],
				},
				include: path.join(__dirname, 'assets'),
				exclude: /(node_modules|bower_components)/,
			},
		],
	},
};

//merge hot reload config

if (DEV_HOT) {
	config.devServer = {
		hot: true,
		inline: true,
		port: 8080,
		proxy: [{
			// for all not hot-update request
			path:    /^(?!.*\.hot-update\.js)(.*)$/,
			target: 'http://localhost:'+ process.env.PORT || 9000,
		}],
		// contentBase:'http://localhost:9000',
		port: process.env.DEV_HOT_PORT || 8090,
		open: true,
		watchOptions: {
			aggregateTimeout: 300,
			poll: 1000
		},
		open: true,
		stats: { colors: true }
	};
	config.plugins.unshift(new webpack.HotModuleReplacementPlugin());
	config.plugins.concat([
		new ProgressBarWebpackPlugin({ clear: false })
	]);

	var babelLoader = config.module.loaders[config.module.loaders.length - 1];
	babelLoader.query.presets.unshift('react-hmre');
	babelLoader.query.plugins = babelLoader.query.plugins || [];
	babelLoader.query.plugins.push([
		'react-transform', {
			transforms: [{
				transform : 'react-transform-hmr',
				imports   : ['react'],
				locals    : ['module']
			}]
		}
	]);
	config.output.publicPath = "http://localhost:" + config.devServer.port + config.output.publicPath
}

module.exports = config;
