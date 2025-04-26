import { MathJaxContext } from "better-react-mathjax";
import Script from "next/script";

const config = {
  "fast-preview": {
    disabled: true,
  },
  tex2jax: {
    inlineMath: [
      ["$", "$"],
      ["\\(", "\\)"],
    ],
    displayMath: [
      ["$$", "$$"],
      ["\\[", "\\]"],
    ],
  },
  messageStyle: "none",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en"  style={{padding: '0px'}} >
      <MathJaxContext config={config}>
      <body style={{margin: '0px'}}>
        {children}
      </body>
      </MathJaxContext>
    </html>
  );
}
