import React from 'react';

export const Card: React.FC<{ children: React.ReactNode; className?: string; glass?: boolean; style?: React.CSSProperties; onClick?: () => void }> = ({ 
  children, 
  className = '',
  glass = false,
  style,
  onClick,
}) => {
  return (
    <div className={`card ${glass ? 'glass-panel' : ''} ${className}`} style={style} onClick={onClick}>
      {children}
    </div>
  );
};
