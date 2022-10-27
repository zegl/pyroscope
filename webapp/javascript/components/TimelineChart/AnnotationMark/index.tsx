/* eslint-disable default-case, consistent-return */
import Color from 'color';
import React, { useState } from 'react';
import classNames from 'classnames/bind';
import AnnotationInfo from '@webapp/pages/continuous/contextMenu/AnnotationInfo';
import useTimeZone from '@webapp/hooks/timeZone.hook';

import styles from './styles.module.scss';

const cx = classNames.bind(styles);

interface IAnnotationMarkProps {
  type: 'message';
  color: Color;
  value: {
    content: string;
    timestamp: number;
  };
}

const getIcon = (type: IAnnotationMarkProps['type']) => {
  switch (type) {
    case 'message':
      return styles.message;
  }
};

const AnnotationMark = ({ type, color, value }: IAnnotationMarkProps) => {
  const { offset } = useTimeZone();
  const [visible, setVisible] = useState(false);
  const [target, setTarget] = useState<Element>();
  const [hovered, setHovered] = useState(false);

  const onClick = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.stopPropagation();
    setTarget(e.target as Element);
    setVisible(true);
  };

  const annotationInfoPopover = target ? (
    <AnnotationInfo
      popoverAnchorPoint={{ x: 0, y: 27 }}
      container={target}
      value={value}
      timezone={offset === 0 ? 'utc' : 'browser'}
      timestamp={value.timestamp}
      isOpen={visible}
      onClose={() => setVisible(false)}
      popoverClassname={styles.form}
    />
  ) : null;

  const onHoverStyle = {
    backgroundColor: hovered ? color.darken(0.2).hex() : color.hex(),
    zIndex: hovered ? 2 : 1,
  };

  return (
    <>
      <div
        data-testid="annotation_mark_wrapper"
        onClick={onClick}
        style={onHoverStyle}
        className={cx(styles.wrapper, getIcon(type))}
        role="none"
        onMouseEnter={() => setHovered(true)}
        onMouseLeave={() => setHovered(false)}
      />
      {annotationInfoPopover}
    </>
  );
};

export default AnnotationMark;
