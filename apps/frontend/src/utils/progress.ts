
const progressBarId = 'global-progress-bar'

export const NProgress = {
    start() {
        let bar = document.getElementById(progressBarId)
        if (!bar) {
            bar = document.createElement('div')
            bar.id = progressBarId
            bar.style.position = 'fixed'
            bar.style.top = '0'
            bar.style.left = '0'
            bar.style.height = '2px'
            bar.style.background = '#409eff' // Element Plus Primary Color
            bar.style.zIndex = '99999'
            bar.style.transition = 'width 0.2s ease'
            bar.style.width = '0%'
            document.body.appendChild(bar)
        }

        // 强制重绘
        bar.offsetHeight

        setTimeout(() => {
            if (bar) bar.style.width = '30%'
        }, 10)
    },

    done() {
        const bar = document.getElementById(progressBarId)
        if (bar) {
            bar.style.width = '100%'
            setTimeout(() => {
                if (bar) {
                    bar.style.opacity = '0'
                    setTimeout(() => {
                        if (bar && bar.parentNode) {
                            bar.parentNode.removeChild(bar)
                        }
                    }, 200)
                }
            }, 200)
        }
    }
}
