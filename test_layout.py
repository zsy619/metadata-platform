from playwright.sync_api import sync_playwright
import time

def test_layout():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # 监听控制台错误
        console_errors = []
        page.on("console", lambda msg: console_errors.append(msg.text) if msg.type == "error" else None)

        try:
            # 访问前端
            print("正在访问前端页面...")
            page.goto('http://localhost:3001/', timeout=30000)
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            # 截图
            page.screenshot(path='/tmp/layout-test.png', full_page=True)
            print("截图已保存到 /tmp/layout-test.png")

            # 检查页面标题
            title = page.title()
            print(f"页面标题: {title}")

            # 检查关键元素
            sidebar = page.locator('.sidebar-container').count()
            print(f"侧边栏组件数量: {sidebar}")

            header = page.locator('.navbar').count()
            print(f"顶部导航栏组件数量: {header}")

            main = page.locator('.app-main').count()
            print(f"主内容区组件数量: {main}")

            # 检查是否有错误
            if console_errors:
                print("\n控制台错误:")
                for err in console_errors:
                    print(f"  - {err}")
            else:
                print("\n无控制台错误")

            print("\n布局组件验证完成！")

        except Exception as e:
            print(f"测试失败: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    test_layout()
