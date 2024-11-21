import json
import re

content = open("call.xml", "r", encoding="utf-8").read()
lines = content.split("\n")

nodes = []


def get_tag_name(text):
    # 使用正则表达式匹配<>之间的内容
    pattern = r'<(.*?)>'

    # 查找所有匹配项
    matches = re.findall(pattern, text)
    if matches:
        return matches[0]
    return ''


n = 0
while n < len(lines):
    line = lines[n]
    if line.startswith("<"):  # 根标签
        node_name = get_tag_name(line)
        if not node_name:
            continue
        #print(node_name)
        node_contents = []
        # 遇到根标签 开始获取全部内容
        while 1:
            line = lines[n]
            if f"</{node_name}>" in line:
                if line.startswith(f"<{node_name}>") and line.endswith(f"</{node_name}>"):  # 无子标签
                    data = line.split(f"<{node_name}>")[1].split(f"</{node_name}>")[0]
                    nodes.append(
                        {node_name: f"<{node_name}>{data}</{node_name}>" }
                    )
                    break

                nodes.append(
                    {node_name: '\n'.join(node_contents), }
                )
                break
            if not line.startswith(f"<{node_name}>"):  # 无子标签
                node_contents.append(line)
            n += 1
    else:
        pass

    n += 1

r = []
for i in nodes:
    if "Termination-Route" in i:
        r.append(i["Termination-Route"])

open("data.txt",'w').write("".join(r))

