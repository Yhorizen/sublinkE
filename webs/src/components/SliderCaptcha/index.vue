<template>
  <div class="slider-captcha">
    <div ref="imgRef" class="slider-img" :style="{ height: imgHeightCss + 'px' }">
      <img
        :src="bgSrc"
        alt=""
        draggable="false"
        :style="{ height: imgHeightCss + 'px' }"
      />
      <!-- 拼块：用户通过下方滑块拖动，对齐背景图中的缺口 -->
      <div
        class="slider-piece"
        :style="{
          width: pieceCss + 'px',
          height: pieceCss + 'px',
          left: pieceLeft + 'px',
          top: slotTopCss + 'px',
        }"
      ></div>
    </div>
    <!-- 自定义滑块轨道：指针事件采集拖动轨迹 -->
    <div
      ref="trackRef"
      class="slider-track"
      :class="{ 'is-disabled': !bgBase64 }"
      @pointerdown="onTrackDown"
      @pointermove="onTrackMove"
      @pointerup="onTrackUp"
      @pointercancel="onTrackUp"
    >
      <div class="slider-fill" :style="{ width: thumbLeft + 'px' }"></div>
      <div class="slider-thumb" :style="{ left: thumbLeft + 'px' }">
        <el-icon><DArrowRight /></el-icon>
      </div>
    </div>
    <div class="slider-foot">
      <span class="slider-hint">{{ hint }}</span>
      <el-icon class="slider-refresh" @click="emit('refresh')"><Refresh /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { DArrowRight, Refresh } from "@element-plus/icons-vue";

const props = withDefaults(
  defineProps<{
    modelValue?: string;
    bgBase64?: string;
    bgWidth?: number;
    bgHeight?: number;
    pieceSize?: number;
  }>(),
  {
    modelValue: "",
    bgBase64: "",
    bgWidth: 320,
    bgHeight: 160,
    pieceSize: 40,
  }
);

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
  (e: "update:trajectory", value: string): void;
  (e: "refresh"): void;
}>();

const { t } = useI18n();
const hint = computed(() => t("login.sliderTip"));

// 兼容无 data URI 前缀的 base64（后端已带前缀，这里做防御性处理）
const bgSrc = computed(() => {
  if (!props.bgBase64) return "";
  return props.bgBase64.startsWith("data:")
    ? props.bgBase64
    : "data:image/png;base64," + props.bgBase64;
});

const imgRef = ref<HTMLElement>();
const trackRef = ref<HTMLElement>();
const displayWidth = ref(0); // 图片容器实际渲染宽度
const thumbLeft = ref(0); // 滑块拇指位置(CSS px)
const dragging = ref(false); // 是否正在拖动
const trajectory = ref<{ x: number; y: number; t: number }[]>([]); // 拖动轨迹点

// 背景图变化(刷新验证码)时重置滑块与轨迹
watch(
  () => props.bgBase64,
  () => {
    thumbLeft.value = 0;
    dragging.value = false;
    trajectory.value = [];
  }
);

// 测量容器宽度，适配响应式布局
const resize = () => {
  if (imgRef.value) {
    displayWidth.value = imgRef.value.clientWidth;
  }
};
onMounted(() => {
  resize();
  if (typeof ResizeObserver !== "undefined") {
    const ro = new ResizeObserver(resize);
    if (imgRef.value) ro.observe(imgRef.value);
  }
});

// 缩放比例：实际渲染宽度 / 图像原始宽度
const scale = computed(() =>
  displayWidth.value > 0 ? displayWidth.value / props.bgWidth : 1
);

// 拼块在容器中的CSS尺寸
const pieceCss = computed(() =>
  Math.max(20, Math.round(props.pieceSize * scale.value))
);

// 拼块CSS位置：跟随滑块比例移动
const pieceLeft = computed(() => {
  const trackW = trackRef.value?.clientWidth ?? displayWidth.value;
  if (trackW <= 0) return 0;
  return (thumbLeft.value / trackW) * (displayWidth.value - pieceCss.value);
});

// 缺口槽位垂直居中，与后端绘制位置一致
const slotTopCss = computed(
  () => ((props.bgHeight - props.pieceSize) / 2) * scale.value
);

const imgHeightCss = computed(() => Math.round(props.bgHeight * scale.value));

// 记录一个轨迹点
function recordPoint(clientX: number, clientY: number) {
  trajectory.value.push({
    x: Math.round(clientX),
    y: Math.round(clientY),
    t: Date.now(),
  });
}

// 按下：开始拖动并记录起点
function onTrackDown(e: PointerEvent) {
  if (!props.bgBase64) return;
  dragging.value = true;
  trajectory.value = [];
  trackRef.value?.setPointerCapture?.(e.pointerId);
  updateThumb(e.clientX);
  recordPoint(e.clientX, e.clientY);
}

// 移动：更新拇指位置并记录轨迹点
function onTrackMove(e: PointerEvent) {
  if (!dragging.value) return;
  updateThumb(e.clientX);
  recordPoint(e.clientX, e.clientY);
}

// 松开：计算最终图像像素X并提交位置与轨迹
function onTrackUp(e: PointerEvent) {
  if (!dragging.value) return;
  dragging.value = false;
  updateThumb(e.clientX);
  recordPoint(e.clientX, e.clientY);

  const trackW = trackRef.value?.clientWidth ?? 0;
  if (trackW > 0) {
    const imgX = Math.round(
      (thumbLeft.value / trackW) * (props.bgWidth - props.pieceSize)
    );
    emit("update:modelValue", String(imgX));
  }
  emit("update:trajectory", JSON.stringify(trajectory.value));
}

// 根据指针X更新拇指位置(限制在轨道内)
function updateThumb(clientX: number) {
  const rect = trackRef.value?.getBoundingClientRect();
  if (!rect || rect.width <= 0) return;
  const x = Math.min(Math.max(clientX - rect.left, 0), rect.width);
  thumbLeft.value = x;
}
</script>

<style scoped>
.slider-captcha {
  width: 100%;
  user-select: none;
}

.slider-img {
  position: relative;
  width: 100%;
  overflow: hidden;
  border-radius: 4px;
  background: #eef1f6;
}

.slider-img img {
  display: block;
  width: 100%;
}

.slider-piece {
  position: absolute;
  box-sizing: border-box;
  background: rgba(110, 122, 148, 0.92);
  border: 1px solid #3c465a;
  border-radius: 3px;
  transition: left 0.05s linear;
  pointer-events: none;
}

.slider-track {
  position: relative;
  height: 34px;
  margin-top: 8px;
  border-radius: 17px;
  background: var(--el-fill-color-light, #f0f2f5);
  border: 1px solid var(--el-border-color, #dcdfe6);
  cursor: pointer;
  touch-action: none;
}

.slider-track.is-disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.slider-fill {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  border-radius: 17px;
  background: rgba(64, 158, 255, 0.15);
  pointer-events: none;
}

.slider-thumb {
  position: absolute;
  top: -1px;
  left: 0;
  width: 36px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 17px;
  background: #fff;
  border: 1px solid var(--el-border-color, #dcdfe6);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  color: var(--el-color-primary, #409eff);
  cursor: grab;
  touch-action: none;
}

.slider-thumb:active {
  cursor: grabbing;
}

.slider-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
}

.slider-hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.slider-refresh {
  cursor: pointer;
  color: var(--el-text-color-secondary);
  font-size: 15px;
}

.slider-refresh:hover {
  color: var(--el-color-primary);
}
</style>
