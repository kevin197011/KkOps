import React from 'react';

interface FooterProps {
  dark?: boolean;
}

const Footer: React.FC<FooterProps> = ({ dark = false }) => {
  return (
    <div
      style={{
        textAlign: 'center',
        padding: '8px 16px',
        background: dark ? '#161b22' : '#f0f2f5',
        color: dark ? '#6e7681' : '#8c8c8c',
        fontSize: '12px',
        borderTop: dark ? '1px solid #30363d' : '1px solid #e8e8e8',
      }}
    >
      本系统由 KK 驱动支持
    </div>
  );
};

export default Footer;
