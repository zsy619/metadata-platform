<template>
    <div class="animated-logo-container">
        <svg width="60" height="60" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
            <defs>
                <!-- 外六边形顺时针流光渐变 -->
                <linearGradient id="outer-gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" stop-color="#ffffff" stop-opacity="0.3" />
                    <stop offset="50%" stop-color="#ffffff" stop-opacity="1" />
                    <stop offset="100%" stop-color="#ffffff" stop-opacity="0.3" />
                    <animateTransform attributeName="gradientTransform" type="rotate" from="0 0.5 0.5" to="360 0.5 0.5" dur="3s" repeatCount="indefinite" />
                </linearGradient>
                <!-- 内六边形逆时针流光渐变 -->
                <linearGradient id="inner-gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" stop-color="#00BCD4" stop-opacity="0.3" />
                    <stop offset="50%" stop-color="#00BCD4" stop-opacity="1" />
                    <stop offset="100%" stop-color="#00BCD4" stop-opacity="0.3" />
                    <!-- 注意：这里用 360 -> 0 实现逆时针视觉 -->
                    <animateTransform attributeName="gradientTransform" type="rotate" from="360 0.5 0.5" to="0 0.5 0.5" dur="3s" repeatCount="indefinite" />
                </linearGradient>
            </defs>
            <!-- 外六边形：固定位置，边框色引用动态渐变 -->
            <!-- 顶点计算：中心(50,50) R=45 -->
            <!-- Points: (50,5) (89,27.5) (89,72.5) (50,95) (11,72.5) (11,27.5) -->
            <polygon points="50,5 89,27.5 89,72.5 50,95 11,72.5 11,27.5" fill="none" stroke="url(#outer-gradient)" stroke-width="3" shape-rendering="geometricPrecision" />
            <!-- 内六边形：固定位置，边框色引用动态渐变 -->
            <!-- 顶点计算：中心(50,50) R=33 -->
            <!-- Points: (50,17) (78.5,33.5) (78.5,66.5) (50,83) (21.5,66.5) (21.5,33.5) -->
            <polygon points="50,17 78.5,33.5 78.5,66.5 50,83 21.5,66.5 21.5,33.5" fill="none" stroke="url(#inner-gradient)" stroke-width="2" shape-rendering="geometricPrecision" />
            <!-- M 形状与装饰点 -->
            <!-- 绝对路径，无 transform -->
            <g>
                <!-- M 的路径定义 -->
                <path id="m-path-fixed" d="M28 72 L28 32 L50 50 L72 32 L72 72" fill="none" stroke="#B2EBF2" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                <!-- 装饰节点 -->
                <!-- 左柱 -->
                <circle cx="28" cy="72" r="2.5" fill="#00BCD4" />
                <circle cx="28" cy="52" r="1.5" fill="#B2EBF2" />
                <circle cx="28" cy="32" r="2.5" fill="#00BCD4" />
                <!-- 中间 -->
                <circle cx="39" cy="41" r="1.5" fill="#B2EBF2" />
                <circle cx="50" cy="50" r="2.5" fill="#00BCD4" />
                <circle cx="61" cy="41" r="1.5" fill="#B2EBF2" />
                <!-- 右柱 -->
                <circle cx="72" cy="32" r="2.5" fill="#00BCD4" />
                <circle cx="72" cy="52" r="1.5" fill="#B2EBF2" />
                <circle cx="72" cy="72" r="2.5" fill="#00BCD4" />
                <!-- 内部连线增加科技感 -->
                <path d="M28 52 L39 41 M61 41 L72 52" stroke="#00BCD4" stroke-width="1" opacity="0.6" />
            </g>
            <!-- 沿着 M 路径运动的红点 -->
            <!-- 严格沿着 id="m-path-fixed" 运动 -->
            <circle r="3" fill="#FF5252">
                <animateMotion dur="4s" repeatCount="indefinite" rotate="auto" keyPoints="0;1;0" keyTimes="0;0.5;1" calcMode="linear">
                    <mpath href="#m-path-fixed" />
                </animateMotion>
            </circle>
        </svg>
    </div>
</template>
<style scoped>
.animated-logo-container {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
}
</style>
