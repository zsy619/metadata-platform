/// <reference types="vite/client" />

declare module 'marked' {
  export function marked(markdown: string, options?: any): string
  export interface MarkedOptions {
    gfm?: boolean
    breaks?: boolean
    async?: boolean
  }
  export function setOptions(options: MarkedOptions): void
}
