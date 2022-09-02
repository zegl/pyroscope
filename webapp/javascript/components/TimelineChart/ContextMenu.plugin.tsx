import React, { useEffect, useState } from 'react';
import * as ReactDOM from 'react-dom';

type ContextType = {
  init: (plot: unknown) => void;
  options: unknown;
  name: string;
  version: string;
};

// TODO(eh-am)
const WRAPPER_ID = 'contextmenu_id';
const injectContextMenu = ($: JQueryStatic) => {
  const parent = $(`#${WRAPPER_ID}`).length
    ? $(`#${WRAPPER_ID}`)
    : $(`<div id="${WRAPPER_ID}" />`);

  const par2 = $(`body`);

  return parent.appendTo(par2);
};

function MyElement() {
  //const [state, setState] = useState(0);
  return <div>hey</div>;
}

(function ($: JQueryStatic) {
  function init(this: ContextType, plot: ShamefulAny) {
    plot.hooks.drawOverlay.push(() => {
      const container = injectContextMenu($);
      ReactDOM.render(MyElement(), container?.[0]);
    });
  }

  ($ as ShamefulAny).plot.plugins.push({
    init,
    options: {},
    name: 'context_menu',
    version: '1.0',
  });
})(jQuery);
