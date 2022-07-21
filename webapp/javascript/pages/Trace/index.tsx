/* eslint-disable no-underscore-dangle */
import * as React from 'react';
import _get from 'lodash/get';
import _memoize from 'lodash/memoize';
import filterSpans from './utils/filter-spans';
import TracePageHeader from './TracePageHeader/TracePageHeader';
import { mockData } from './mockTrace';
import {
  TUpdateViewRangeTimeFunction,
  IViewRange,
  ViewRangeTimeUpdate,
  ETraceViewType,
  FetchedTrace,
  ReduxState,
  TNil,
} from './types';
import TraceTimelineViewer from './TraceTimelineViewer';
import ScrollManager from './ScrollManager';
import { cancel as cancelScroll, scrollBy, scrollTo } from './scroll-page';

type TProps = TDispatchProps & TOwnProps & TReduxProps;

type TState = {
  headerHeight: number | TNil;
  slimView: boolean;
  viewType: ETraceViewType;
  viewRange: IViewRange;
};

export default class TracePageImpl extends React.PureComponent<TProps, TState> {
  state: TState;

  _headerElm: HTMLElement | TNil;

  filterSpans: typeof filterSpans;

  _searchBar: React.RefObject<Input>;

  _scrollManager: ScrollManager;

  traceDagEV: ShamefulAny | TNil;

  constructor(props) {
    super(props);
    const { embedded } = props;

    this.state = {
      headerHeight: null,
      slimView: Boolean(embedded && embedded.timeline.collapseTitle),
      viewType: ETraceViewType.TraceTimelineViewer,
      viewRange: {
        time: {
          current: [0, 1],
        },
      },
    };

    this.filterSpans = _memoize(
      filterSpans,
      // Do not use the memo if the filter text or trace has changed.
      // trace.data.spans is populated after the initial render via mutation.
      (textFilter) =>
        `${textFilter} ${_get(this.props.trace, 'traceID')} ${_get(
          this.props.trace,
          'data.spans.length'
        )}`
    );
    this._scrollManager = new ScrollManager(mockData && mockData.data, {
      scrollBy,
      scrollTo,
    });
  }

  setHeaderHeight = (elm: HTMLElement | TNil) => {
    this._headerElm = elm;
    if (elm) {
      if (this.state.headerHeight !== elm.clientHeight) {
        this.setState({ headerHeight: elm.clientHeight });
      }
    } else if (this.state.headerHeight) {
      this.setState({ headerHeight: null });
    }
  };

  toggleSlimView = () => {
    const { slimView } = this.state;
    this.setState({ slimView: !slimView });
  };

  updateViewRangeTime: TUpdateViewRangeTimeFunction = (
    start: number,
    end: number,
    trackSrc?: string
  ) => {
    const current: [number, number] = [start, end];
    const time = { current };
    this.setState((state: TState) => ({
      viewRange: { ...state.viewRange, time },
    }));
  };

  updateNextViewRangeTime = (update: ViewRangeTimeUpdate) => {
    this.setState((state: TState) => {
      const time = { ...state.viewRange.time, ...update };
      return { viewRange: { ...state.viewRange, time } };
    });
  };

  render() {
    const {
      archiveEnabled,
      archiveTraceState,
      embedded,
      id,
      uiFind,
      trace,
      // location: { state: locationState },
    } = this.props;
    const { slimView, viewType, headerHeight, viewRange } = this.state;
    // const { data } = trace;

    let findCount = 0;
    let graphFindMatches: Set<string> | null | undefined;
    let spanFindMatches: Set<string> | null | undefined;
    if (uiFind) {
      if (viewType === ETraceViewType.TraceGraph) {
        graphFindMatches = getUiFindVertexKeys(
          uiFind,
          _get(this.traceDagEV, 'vertices', [])
        );
        findCount = graphFindMatches ? graphFindMatches.size : 0;
      } else {
        spanFindMatches = this.filterSpans(uiFind, _get(trace, 'data.spans'));
        findCount = spanFindMatches ? spanFindMatches.size : 0;
      }
    }

    const isEmbedded = Boolean(embedded);
    const headerProps = {
      focusUiFindMatches: () => {}, // this.focusUiFindMatches,
      slimView,
      textFilter: uiFind,
      viewType,
      viewRange,
      canCollapse:
        !embedded ||
        !embedded.timeline.hideSummary ||
        !embedded.timeline.hideMinimap,
      clearSearch: () => {}, // this.clearSearch,
      hideMap: Boolean(
        viewType !== ETraceViewType.TraceTimelineViewer ||
          (embedded && embedded.timeline.hideMinimap)
      ),
      hideSummary: Boolean(embedded && embedded.timeline.hideSummary),
      linkToStandalone: '',
      // nextResult: this.nextResult,
      // onArchiveClicked: this.archiveTrace,
      onSlimViewClicked: this.toggleSlimView,
      // onTraceViewChange: this.setTraceView,
      // prevResult: this.prevResult,
      ref: this._searchBar,
      resultCount: findCount,
      showArchiveButton: !isEmbedded && archiveEnabled,
      showShortcutsHelp: !isEmbedded,
      showStandaloneLink: isEmbedded,
      showViewOptions: !isEmbedded,
      toSearch: null,
      trace: mockData.data,
      updateNextViewRangeTime: this.updateNextViewRangeTime,
      updateViewRangeTime: this.updateViewRangeTime,
    };

    const view = (
      <TraceTimelineViewer
        registerAccessors={this._scrollManager.setAccessors}
        scrollToFirstVisibleSpan={this._scrollManager.scrollToFirstVisibleSpan}
        findMatchesIDs={spanFindMatches}
        trace={mockData.data}
        updateNextViewRangeTime={this.updateNextViewRangeTime}
        updateViewRangeTime={this.updateViewRangeTime}
        viewRange={viewRange}
      />
    );

    return (
      <>
        <div className="Tracepage--headerSection" ref={this.setHeaderHeight}>
          <TracePageHeader {...headerProps} />
        </div>
        {headerHeight ? <section>{view}</section> : null}
      </>
    );
  }
}
