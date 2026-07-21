import React from 'react';

export const Card: React.FC<{ children: React.ReactNode; className?: string; glass?: boolean; style?: React.CSSProperties }> = ({ 
  children, 
  className = '',
  glass = false,
  style
}) => {
  return (
    <div className={`card ${glass ? 'glass-panel' : ''} ${className}`} style={style}>
      {children}
    </div>
  );
};
