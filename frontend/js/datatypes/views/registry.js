/*
Views allow ConnectorDB to display plots and other visualizations of data. Each view
registers itself using addView from ../datatypes.js

This file is imported at the beginning of the app, and it imports all of the available views
so that they can be registered.
*/

import "./TableView/index";
import "./InfoView";
import "./BoolView";
import "./LineView";
import "./HistogramView";
import "./BarView";
import "./SingleObjectView";
import "./SentimentLineChart";
import "./ScatterView";
import "./Map";
import "./WebHistoryView";
import "./HeatmapView";
